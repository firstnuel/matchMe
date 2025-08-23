

export interface MatchError {
    error: string;
    details: string;
}

export interface RecommendationsResponse {
  message: string;
  recommendations: string[]
}

export interface DistanceResponse {
  distance: number;
  unit: "km";
  current_user: string;
  target_user: string;
}
