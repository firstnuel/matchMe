import { useUserProfileBio } from "../hooks/useCurrentUser";
import "../styles.css";
import Section from "./Section";
import FieldDisplay from "./FieldDisplay";
import ImageCarousel from "../../../shared/components/ImageCarousel";
import { getInitials } from "../../../shared/utils/utils";
import { useCityFromCoordinates } from "../hooks/useCityFromCoordinates";
import { useNavigate, useParams } from "react-router-dom";
import IsLoading from "../../../shared/components/IsLoading";
import { useUIStore } from "../../../shared/hooks/uiStore";
import { Icon } from "@iconify/react/dist/iconify.js";

const ViewOtherProfile = () => {
    const { id } = useParams();
    const { data: currentUser, isPending } = useUserProfileBio(id ?? "");
    const navigate = useNavigate();
    const user = currentUser && 'user' in currentUser ? currentUser.user : undefined;
    const { locationDisplay, fetchingCity } = useCityFromCoordinates(user?.coordinates);
    const { view } = useUIStore()
    const initials = getInitials(user?.first_name ?? '', user?.last_name ?? '');
    const formatGender = (gender?: string) => {
        if (!gender) return 'Not specified';
        return gender.charAt(0).toUpperCase() + gender.slice(1).replace('_', ' ');
    };

    if (isPending) {
        return <IsLoading message="Loading profile..." />;
    }

    return(
        <div className="view-content">
            {view === "chat" &&
                <button  className="back-btn"
                onClick={() => navigate("/chat")}>
                    <Icon icon="mdi:arrow-back" className="back-icon" />
                    Back
                </button>
            }
            <div className="card-image crd-img">
                {(Array.isArray(user?.photos) && user.photos.length > 0) || user?.profile_photo ? (
                <ImageCarousel 
                    photos={user?.photos ?? []}
                    fallbackPhoto={user?.profile_photo ?? undefined}
                    altText={`${user?.first_name || 'User'} profile`}
                />
                ) : (
                <div className="profile-placeholder">
                    <div className="placeholder-avatar">
                    {initials}
                    </div>
                </div>
                )}
            </div>
            
            <Section title="Basic Information" subtitle="Core profile details">
                <FieldDisplay label="First Name" value={user?.first_name} />
                <FieldDisplay label="Last Name" value={user?.last_name} />
                <FieldDisplay label="Age" value={user?.age} />
                <FieldDisplay label="Gender" value={formatGender(user?.gender)} />
            </Section>

            <Section title="About" subtitle="A short bio">
                <FieldDisplay label="About Me" value={user?.about_me} />
            </Section>

            <Section title="Location" subtitle="Where this user is currently based">
                <FieldDisplay
                    label="Current Location"
                    value={fetchingCity ? "Fetching Location" : locationDisplay}
                />
            </Section>

            <Section title="Preferences" subtitle="What this user is looking for">
                <FieldDisplay label="Looking For" value={user?.looking_for} type="array" />
                <FieldDisplay label="Interests" value={user?.interests} type="array" />
                <FieldDisplay label="Music Preferences" value={user?.music_preferences} type="array" />
                <FieldDisplay label="Food Preferences" value={user?.food_preferences} type="array" />
                <FieldDisplay label="Communication Style" value={user?.communication_style} />
            </Section>

            <Section title="Profile Prompts" subtitle="Get to know this user better" className="last-el">
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

export default ViewOtherProfile