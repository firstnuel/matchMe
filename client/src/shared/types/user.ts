export interface Point {
  longitude: number;
  latitude: number;
  city?: string
}

export interface Prompt {
  id?: string;
  question: string;
  answer: string;
}

export interface UserPhoto {
  id: string;
  photo_url: string;
  order: number;
}

export interface User {
  id: string;
  email?: string;
  first_name: string;
  last_name: string;
  created_at?: string | null;
  updated_at?: string | null;
  age?: number;
  about_me?: string | null;
  preferred_age_min?: number | null;
  preferred_distance?: number | null;
  preferred_age_max?: number | null;
  profile_completion: number;
  gender?: string;
  preferred_gender?: string;
  coordinates?: Point | null;
  looking_for?: string[];
  interests?: string[];
  music_preferences?: string[];
  food_preferences?: string[];
  communication_style?: string | null;
  prompts?: Prompt[];
  photos?: UserPhoto[];
  profile_photo?: string | null;
}