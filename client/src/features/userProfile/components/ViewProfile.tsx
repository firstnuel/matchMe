import ProfileImageCard from "./ProfileImageCard"
import { useCurrentUser } from "../hooks/useCurrentUser";
import { useNavigate } from "react-router";
import "../styles.css";
import Section from "./Section";
import FieldDisplay from "./FieldDisplay";

import { useCityFromCoordinates } from "../hooks/useCityFromCoordinates";

const ViewProfile = () => {
    const { data: currentUser } = useCurrentUser();
    const user = currentUser && 'user' in currentUser ? currentUser.user : undefined;
    const navigate = useNavigate();
    const { locationDisplay, fetchingCity } = useCityFromCoordinates(user?.coordinates);

    // Format gender display
    const formatGender = (gender?: string) => {
        if (!gender) return 'Not specified';
        return gender.charAt(0).toUpperCase() + gender.slice(1).replace('_', ' ');
    };

    return(
        <div className="view-content">
            <ProfileImageCard 
                user={user} 
                onEditClick={() => navigate("/edit-profile")}
            />
            
            <Section title="Basic Information" subtitle="Your core profile details">
                <FieldDisplay label="First Name" value={user?.first_name} />
                <FieldDisplay label="Last Name" value={user?.last_name} />
                <FieldDisplay label="Email" value={user?.email} />
                <FieldDisplay label="Age" value={user?.age} />
                <FieldDisplay label="Gender" value={formatGender(user?.gender)} />
            </Section>

            <Section title="About Me" subtitle="Tell others about yourself">
                <FieldDisplay label="About Me" value={user?.about_me} />
            </Section>

            <Section title="Location" subtitle="Help others find you nearby">
                <FieldDisplay label="Current Location" value={fetchingCity? "Fetching Location": locationDisplay} />
            </Section>

            <Section title="Your Preferences" subtitle="What you're looking for">
                <FieldDisplay label="Looking For" value={user?.looking_for} type="array" />
                <FieldDisplay label="Interests" value={user?.interests} type="array" />
                <FieldDisplay label="Music Preferences" value={user?.music_preferences} type="array" />
                <FieldDisplay label="Food Preferences" value={user?.food_preferences} type="array" />
                <FieldDisplay label="Communication Style" value={user?.communication_style} />
            </Section>

            <Section title="Dating Preferences" subtitle="Help us find your perfect matches">
                <FieldDisplay 
                    label="Preferred Age Range" 
                    value={
                        user?.preferred_age_min && user?.preferred_age_max
                            ? `${user.preferred_age_min} - ${user.preferred_age_max}`
                            : undefined
                    }
                    type="range" 
                />
                <FieldDisplay label="Maximum Distance" value={user?.preferred_distance ? `${user.preferred_distance} km(s)` : undefined} />
                <FieldDisplay label="Preferred Gender" value={formatGender(user?.preferred_gender)} />
            </Section>

            <Section title="Profile Prompts" subtitle="Your personality prompts" className="last-el">
                {user?.prompts && user.prompts.length > 0 ? (
                    user.prompts.map((prompt, index) => (
                        <div key={index} className="sec-show">
                            <div className="sec-name">{prompt.question}</div>
                            <div className="sec-value">{prompt.answer}</div>
                        </div>
                    ))
                ) : (
                    <FieldDisplay label="Profile Prompts" value="No prompts added yet" />
                )}
            </Section>
        </div>
    )
}

export default ViewProfile