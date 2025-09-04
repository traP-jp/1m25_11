import type { components } from '~~/shared/types/generated';

// export type StampSummary = paths['/stamps']['get']['responses'][200]['content']['application/json'][0];

// export type StampDetail = paths['/stamps/{stampId}']['get']['responses']['200']['content']['application/json'];

export type Schemas = {
  [K in keyof components['schemas']]: components['schemas'][K]
};
// export interface stampSummary {
//   id: string;
//   name: string;
//   file_id: string;
// }

// export interface tagSummary {
//   id: string;
//   name: string;
// }

// export interface stampDescription {
//   stamp_id: string;
//   description: string;
//   creator_id: string;
//   created_at: string;
//   updated_at: string;
// }

// export interface detailResponse {
//   id: string;
//   name: string;
//   file_id: string;
//   creator_id: string;
//   is_unicode: boolean;
//   created_at: string;
//   updated_at: string; // time.Time (ISO 8601 format)
//   count_monthly: number;
//   count_total: number;

//   descriptions: stampDescription[];
//   tags: tagSummary[];
// }

// // tag の定義は作ってない
