import { type User } from "../../../shared/types/user";

export const Interests: string[] = [
  "travel", "music", "movies", "books", "cooking", "fitness", "art", "photography",
  "gaming", "sports", "hiking", "dancing", "yoga", "meditation", "technology",
  "fashion", "food", "wine", "coffee", "pets", "nature", "adventure", "reading",
];

export const MusicPreferences: string[] = [
  "pop", "rock", "jazz", "classical", "hip-hop", "electronic", "country", "folk",
  "blues", "reggae", "indie", "alternative", "r&b", "soul", "funk", "punk",
  "metal", "latin", "world", "ambient", "afrobeats", "amapiano",
];

export const FoodPreferences: string[] = [
  "vegetarian", "vegan", "italian", "chinese", "japanese", "mexican", "indian",
  "thai", "french", "mediterranean", "american", "korean", "vietnamese",
  "middle-eastern", "african", "fusion", "seafood", "bbq", "desserts", "street-food",
];

export const promptQuestions = [
  "What's your ideal Sunday like?",
  "Two truths and a lie",
  "My biggest fear is...",
  "What's your go-to comfort food?",
  "If you could travel anywhere, where would you go?",
  "What's a hobby you want to pick up?",
  "Your favorite book or movie is...",
  "What's a fun fact about you?",
  "Your dream dinner guest would be...",
  "What's your hidden talent?"
];

export const CommunicationStyles: string[] = [
  "direct", "thoughtful", "humorous", "analytical", "creative", "empathetic",
  "casual", "formal", "energetic", "calm",
];

export interface UserError {
    error: string;
    details: string;
}

export  interface UserResponse {
  message: string;
  user:  User | Partial<User> | null
}

export interface UpdateUserRequest {
  first_name?: string; 
  last_name?: string;
  age?: number;
  gender?: "male" | "female" | "non_binary" | "prefer_not_to_say";
  location?: Location;
  about_me?: string;
  bio?: UserBio;
  preferred_age_min?: number;
  preferred_age_max?: number;
  preferred_gender?: "male" | "female" | "non_binary" | "all";
  preferred_distance?: number;
}

interface Location {
  latitude: number; 
  longitude: number;
}

export interface UserBio {
  looking_for?: string[]; 
  interests?: string[]; 
  music_preferences?: string[]; 
  food_preferences?: string[]; 
  communication_style?:string; 
  prompts?: Prompt[];
}

interface Prompt {
  question: string;
  answer: string; 
}

export interface PhotoUploadResponse {
  message: string;
  photos: UserPhoto[];
  count: number;
}

export interface UserPhoto {
  id: string;
  photo_url: string;
  order: number;
}