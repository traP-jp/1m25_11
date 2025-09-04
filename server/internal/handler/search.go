package handler

import (
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type searchStampsParams struct {
	Q                  *string  `query:"q"`
	Name               *string  `query:"name"`
	Tag                []string `query:"tag"`
	Description        *string  `query:"description"`
	CreatedSince       *string  `query:"created_since"`
	CreatedUntil       *string  `query:"created_until"`
	UpdatedSince       *string  `query:"updated_since"`
	UpdatedUntil       *string  `query:"updated_until"`
	StampTypeUnicode   *string  `query:"stamp_type_unicode"`
	StampTypeAnimation *string  `query:"stamp_type_animation"`
	CountMonthlyMin    *int     `query:"count_monthly_min"`
	CountMonthlyMax    *int     `query:"count_monthly_max"`
	SortBy             *string  `query:"sortby"`
}

type searchResultResponse struct {
	Stamps []stampSummaryResponse `json:"stamps"`
}

type stampSummaryResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	FileID string `json:"file_id"`
}

type scoredStamp struct {
	Stamp repository.StampForSearch
	Score float64
}

func (h *Handler) SearchStamps(c echo.Context) error {
	var params searchStampsParams
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request parameters: "+err.Error())
	}
	repoParams := repository.SearchStampsParams{}
	if params.Q != nil {
		repoParams.Query = *params.Q
	}
	if params.Name != nil {
		repoParams.Name = *params.Name
	}
	repoParams.Tags = params.Tag
	if params.Description != nil {
		repoParams.Description = *params.Description
	}

	const layout = "2006-01-02"
	if params.CreatedSince != nil {
		if t, err := time.Parse(layout, *params.CreatedSince); err == nil {
			repoParams.CreatedSince = &t
		}
	}
	if params.CreatedUntil != nil {
		if t, err := time.Parse(layout, *params.CreatedUntil); err == nil {
			repoParams.CreatedUntil = &t
		}
	}
	if params.UpdatedSince != nil {
		if t, err := time.Parse(layout, *params.UpdatedSince); err == nil {
			repoParams.UpdatedSince = &t
		}
	}
	if params.UpdatedUntil != nil {
		if t, err := time.Parse(layout, *params.UpdatedUntil); err == nil {
			repoParams.UpdatedUntil = &t
		}
	}

	if params.StampTypeUnicode != nil {
		repoParams.StampTypeUnicode = *params.StampTypeUnicode
	}
	if params.StampTypeAnimation != nil {
		repoParams.StampTypeAnimation = *params.StampTypeAnimation
	}
	if params.CountMonthlyMin != nil {
		repoParams.CountMonthlyMin = params.CountMonthlyMin
	}
	if params.CountMonthlyMax != nil {
		repoParams.CountMonthlyMax = params.CountMonthlyMax
	}
	if params.SortBy != nil {
		repoParams.SortBy = *params.SortBy
	}

	foundStamps, err := h.repo.SearchStamps(c.Request().Context(), repoParams)
	if err != nil {
		log.Printf("error in SearchStamps repository call: %v", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to search stamps")
	}

	stampsRes := make([]stampSummaryResponse, len(foundStamps))

	if repoParams.SortBy == "relativity" || repoParams.SortBy == "" {
		scoredStamps := make([]scoredStamp, len(foundStamps))
		for i, stamp := range foundStamps {
			score := calculateRelativityScore(stamp, repoParams)
			scoredStamps[i] = scoredStamp{Stamp: stamp, Score: score}
		}

		sort.Slice(scoredStamps, func(i, j int) bool {
			if scoredStamps[i].Score != scoredStamps[j].Score {
				return scoredStamps[i].Score > scoredStamps[j].Score
			}
			return scoredStamps[i].Stamp.Name < scoredStamps[j].Stamp.Name
		})

		for i, ss := range scoredStamps {
			stampsRes[i] = stampSummaryResponse{
				ID:     ss.Stamp.ID.String(),
				Name:   ss.Stamp.Name,
				FileID: ss.Stamp.FileID.String(),
			}
		}
	} else {
		for i, s := range foundStamps {
			stampsRes[i] = stampSummaryResponse{
				ID:     s.ID.String(),
				Name:   s.Name,
				FileID: s.FileID.String(),
			}
		}
	}

	response := searchResultResponse{
		Stamps: stampsRes,
	}

	return c.JSON(http.StatusOK, response)
}

func calculateRelativityScore(stamp repository.StampForSearch, params repository.SearchStampsParams) float64 {
	divisor := 0.0    
	totalScore := 0.0 
	calculateSubScore := func(query, targetText string) float64 {
		terms := strings.Fields(query)
		if len(terms) == 0 {
			return 0
		}
		var sumOfX float64
		for _, term := range terms {
			count := strings.Count(strings.ToLower(targetText), strings.ToLower(term))
			x := 1.0 - math.Exp(float64(-count))
			sumOfX += x
		}
		return sumOfX / float64(len(terms))
	}
	if params.Name != "" {
		scoreName := calculateSubScore(params.Name, stamp.Name)
		totalScore += scoreName
		divisor++
	}
	if params.Description != "" {
		scoreDescription := calculateSubScore(params.Description, stamp.Descriptions)
		totalScore += scoreDescription
		divisor++
	}
	if len(params.Tags) > 0 {
		scoreTag := calculateSubScore(strings.Join(params.Tags, " "), stamp.Tags)
		totalScore += scoreTag
		divisor++
	}
	var scoreQ float64
	qTerms := strings.Fields(params.Query)
	if len(qTerms) > 0 {
		var sumOfXi float64
		for _, term := range qTerms {
			xName := 1.0 - math.Exp(float64(-strings.Count(strings.ToLower(stamp.Name), strings.ToLower(term))))
			xTag := 1.0 - math.Exp(float64(-strings.Count(strings.ToLower(stamp.Tags), strings.ToLower(term))))
			xDesc := 1.0 - math.Exp(float64(-strings.Count(strings.ToLower(stamp.Descriptions), strings.ToLower(term))))
			xi := (xName + xTag + xDesc) / 3.0
			sumOfXi += xi
		}
		scoreQ = sumOfXi / float64(len(qTerms))
		totalScore += scoreQ
		divisor++
	}

	finalScore := totalScore / divisor
	return finalScore
}
