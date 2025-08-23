import PhotoUploadSection from "./PhotoUploadCard";
import FormField from "./FormField";
import Location from "./Location";
import { type Point, type Prompt, type User } from "../../../shared/types/user";
import { useState } from "react";
import TagSelector from "./TagSelector";
import Section from "./Section";
import { useUserProfile } from "../hooks/useUserProfile";
import { useField } from "../../../shared/hooks/useField";
import { useCurrentUser } from "../hooks/useCurrentUser";
import { useSelect } from "../../../shared/hooks/useSelect";
import {
  Interests,
  MusicPreferences,
  FoodPreferences,
  CommunicationStyles,
} from "../types/user";
import "../styles.css";
import ProfilePrompts from "./ProfilePrompts";

const EditProfile = () => {
  const { data: currentUser } = useCurrentUser();
  const user = currentUser && 'user' in currentUser ? currentUser.user : undefined;
  const { setPendingUpdate } = useUserProfile()
  const [ageRange, setAgeRange] = useState<[number, number]>(
    user?.preferred_age_min === 0 || user?.preferred_age_max === 0
      ? [22, 35]
      : [user?.preferred_age_min ?? 22, user?.preferred_age_max ?? 35]
  );
  const [maxDistance, setMaxDistance] = useState<number>(
    user?.preferred_distance === 0 ? 25 : user?.preferred_distance ?? 25
  );
  const firstName = useField("firstName", "text", user?.first_name)
  const lastName = useField("lastName", "text", user?.last_name)
  const age = useField("age", "number", user?.age)
  const aboutMe = useField('bio', 'textarea', user?.about_me?? "")
  const genderField = useSelect(user?.gender?? "male")
  const preferredGenderField = useSelect(user?.preferred_gender?? "all")
  const [selectedInterests, setSelectedInterests] = useState<string[]>([])
  const [selectedLookingFor, setSelectedLookingFor] = useState<string[]>([])
  const [selectedMusicPreferences, setSelectedMusicPreferences] = useState<string[]>([])
  const [selectedFoodPreferences, setSelectedFoodPreferences] = useState<string[]>([])
  const [selectedCommunicationStyle, setSelectedCommunicationStyle] = useState<string[]>([])
  const [location, setLocation] = useState<Point>(user?.coordinates ?? { longitude: 0, latitude: 0, city: "" })
  const [prompts, setPrompts] = useState<Prompt[]>(user?.prompts ?? []);

  // Build userUpdate data with non-empty field validation
  const buildUserUpdateData = () => {
    const updateData: Partial<User> = {};

    // Basic information
    if (firstName.value && (firstName.value as string).trim()) {
      updateData.first_name = (firstName.value as string).trim();
    }
    
    if (lastName.value && (lastName.value as string).trim()) {
      updateData.last_name = (lastName.value as string).trim();
    }
    
    if (age.value && Number(age.value) > 0) {
      updateData.age = Number(age.value);
    }
    
    if (genderField.value) {
      updateData.gender = genderField.value;
    }

    // About me
    if (aboutMe.value && (aboutMe.value as string).trim()) {
      updateData.about_me = (aboutMe.value as string).trim();
    }

    // Location - check if coordinates are valid (not 0,0)
    if (location && (location.longitude !== 0 || location.latitude !== 0)) {
      updateData.coordinates = location;
    }

    // Preferences arrays - only add if not empty
    if (selectedLookingFor.length > 0) {
      updateData.looking_for = selectedLookingFor;
    }
    
    if (selectedInterests.length > 0) {
      updateData.interests = selectedInterests;
    }
    
    if (selectedMusicPreferences.length > 0) {
      updateData.music_preferences = selectedMusicPreferences;
    }
    
    if (selectedFoodPreferences.length > 0) {
      updateData.food_preferences = selectedFoodPreferences;
    }
    
    if (selectedCommunicationStyle.length > 0) {
      updateData.communication_style = selectedCommunicationStyle[0]; // Only one allowed
    }

    // Dating preferences
    if (ageRange[0] !== 0 && ageRange[1] !== 0) {
      updateData.preferred_age_min = ageRange[0];
      updateData.preferred_age_max = ageRange[1];
    }
    
    if (maxDistance > 0) {
      updateData.preferred_distance = maxDistance;
    }
    
    if (preferredGenderField.value) {
      updateData.preferred_gender = preferredGenderField.value;
    }

    // Prompts - filter out empty prompts
    const validPrompts = prompts.filter(prompt => 
      prompt.question && prompt.question.trim() && 
      prompt.answer && prompt.answer.trim()
    );
    if (validPrompts.length > 0) {
      updateData.prompts = validPrompts;
    }

    return updateData;
  };

    const savePendingUpdate = () => {
    const userData = buildUserUpdateData();
    setPendingUpdate(userData);
    console.log('Saved to pending update:', userData);
    };


  return (
    <div className="edit-content">

      <Section title="Basic Information" subtitle="Your core profile details">
        <FormField label="First Name" type="text" value={firstName.value} onChange={firstName.onChange} maxLength={50} />
        <FormField label="Last Name" type="text" value={lastName.value} onChange={lastName.onChange} maxLength={30} />
        <FormField label="Age" type="number" value={age.value} onChange={age.onChange} min={18} max={100} />
        <FormField
          label="Gender"
          type="select"
          value={genderField.value}
          onChange={genderField.onChange}
          options={[
            { value: "male", label: "Male" },
            { value: "female", label: "Female" },
            { value: "non_binary", label: "Non-binary" },
            { value: "prefer_not_to_say", label: "Prefer not to say" },
          ]}
        />
      </Section>
      <Section title="Photos" subtitle="Upload up to 5 photos. The first photo will be your main photo — drag to reorder.">
        <PhotoUploadSection />
      </Section>

      <Section title="About Me" subtitle="Tell others about yourself">
        <FormField
          label="About Me"
          type="textarea"
          value={aboutMe.value}
          onChange={aboutMe.onChange}
          maxLength={500}
        />
      </Section>

      <Section title="Location" subtitle="Help others find you nearby">
        <Location location={location} setLocation={setLocation} />
      </Section>

      <Section title="Your Preferences" subtitle="Select all that apply">
        <TagSelector
          options={["relationship", "friendship", "casual", "networking"]}
          maxSelectable={2}
          label="Looking For"
          selectedTags={selectedLookingFor}
          setSelectedTags={setSelectedLookingFor}
        />
        <TagSelector
          options={Interests}
          maxSelectable={5}
          selectedTags={selectedInterests}
          setSelectedTags={setSelectedInterests}
          label="Interests (Max 5)"
        />
        <TagSelector
          options={MusicPreferences}
          maxSelectable={5}
          selectedTags={selectedMusicPreferences}
          setSelectedTags={setSelectedMusicPreferences}
          label="Music Preferences (Max 5)"
        />
        <TagSelector
          options={FoodPreferences}
          maxSelectable={5}
          selectedTags={selectedFoodPreferences}
          setSelectedTags={setSelectedFoodPreferences}
          label="Food Preferences (Max 5)"
        />
        <TagSelector
          options={CommunicationStyles}
          maxSelectable={1}
          selectedTags={selectedCommunicationStyle}
          setSelectedTags={setSelectedCommunicationStyle}
          label="Communication Style (Max 1)"
        />
      </Section>

      <Section title="Dating Preferences" subtitle="Help us find your perfect matches">
        <FormField
          label="Preferred Age Range"
          type="range"
          value={ageRange}
          min={18}
          max={100}
          onChange={(e) => setAgeRange(e.target.value as [number, number])}
        />
        <FormField
          label="Maximum Distance"
          type="slider"
          value={maxDistance}
          min={1}
          max={100}
          unit="km(s)"
          onChange={(e) => setMaxDistance(Number((e.target as HTMLInputElement).value))}
        />
        <FormField
          label="Preferred Gender"
          type="select"
          value={preferredGenderField.value}
          onChange={preferredGenderField.onChange}
          options={[
            { value: "male", label: "Male" },
            { value: "female", label: "Female" },
            { value: "non_binary", label: "Non-binary" },
            { value: "all", label: "No Preference" },
          ]}
        />
      </Section>

      <Section title="Profile Prompts" subtitle="Answer 3-5 prompts to show your personality">    
        <ProfilePrompts
          prompts={prompts}
          setPrompts={setPrompts}
        />
      </Section>

    </div>
  );
};

export default EditProfile;