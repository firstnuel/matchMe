import { useState, useEffect } from 'react';
import ProfileCard from './ProfileCard';
import { useUserRecommendations, useUserProfile, useUserDistance } from '../hooks/useMatch';
import { useCurrentUser } from '../../userProfile/hooks/useCurrentUser';
import { useNavigate } from 'react-router';
import '../styles.css';

const Home = () => {
    const { data: currentUserData } = useCurrentUser();
    const { data: recommendationsData, isLoading: isLoadingRecommendations } = useUserRecommendations();
    const [actionedProfiles, setActionedProfiles] = useState<Set<string>>(new Set());
    const [currentProfileIndex, setCurrentProfileIndex] = useState(0);
    const navigate = useNavigate()
    
    // Get current user and profile completion
    const loggedInUser = currentUserData && 'user' in currentUserData ? currentUserData.user : null;
    const profileCompletion = loggedInUser?.profile_completion ?? 0;
    const isProfileIncomplete = profileCompletion <= 90;
    
    // Get recommendations array from server response
    const recommendations = recommendationsData && 'recommendations' in recommendationsData 
        ? recommendationsData.recommendations 
        : [];
    
    // Filter out actioned profiles and get current profile ID
    const availableProfiles = recommendations.filter(id => !actionedProfiles.has(id));
    const currentProfileId = availableProfiles[currentProfileIndex];
    
    // Fetch current profile data only if there's a valid ID
    const { data: currentProfileData, isLoading: isLoadingProfile } = useUserProfile(currentProfileId || '');
    const { data: distanceData } = useUserDistance(currentProfileId || '');
    
    // Get user from profile response
    const currentUser = currentProfileData && 'user' in currentProfileData ? currentProfileData.user : null;
    const distance = distanceData && 'distance' in distanceData ? distanceData.distance : null;

    // Reset to first profile when recommendations change
    useEffect(() => {
        if (recommendations.length > 0) {
            setCurrentProfileIndex(0);
        }
    }, [recommendations.length]);

    const handleAction = (profileId: string, action: 'like' | 'reject') => {
        // Add to actioned profiles
        setActionedProfiles(prev => new Set([...prev, profileId]));
        
        // Move to next profile if available
        if (currentProfileIndex < availableProfiles.length - 1) {
            setCurrentProfileIndex(prev => prev + 1);
        }
        
        console.log(`${action} profile:`, profileId);
    };

    const handleLike = () => {
        if (currentProfileId) {
            handleAction(currentProfileId, 'like');
        }
    };

    const handleReject = () => {
        if (currentProfileId) {
            handleAction(currentProfileId, 'reject');
        }
    };

    // Profile completion check - show message if profile is incomplete
    if (loggedInUser && isProfileIncomplete) {
        return (
            <div className="match-content">
                <div className="card-stack">
                    <div className="message-container">
                        <h2 className="message-title">Complete Your Profile</h2>
                        <p className="message-text">
                            Your profile is {Math.round(profileCompletion)}% complete. Complete at least 90% of your profile to start seeing recommendations.
                        </p>
                        <button 
                            onClick={() => navigate("/edit-profile")}
                            className="complete-profile-btn"
                        >
                            Complete Profile
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    // Loading state
    if (isLoadingRecommendations) {
        return (
            <div className="match-content">
                <div className="card-stack">
                    <div className="loading-container">
                        Loading recommendations...
                    </div>
                </div>
            </div>
        );
    }

    // No more profiles available
    if (availableProfiles.length === 0) {
        return (
            <div className="match-content">
                <div className="card-stack">
                    <div className="message-container">
                        <h2 className="message-title">No More Profiles</h2>
                        <p className="no-profiles-text">
                            You've seen all available matches! Check back later for new recommendations.
                        </p>
                    </div>
                </div>
            </div>
        );
    }

    // Loading current profile
    if (isLoadingProfile && currentUser === null) {
        return (
            <div className="match-content">
                <div className="card-stack">
                    <div className="loading-container">
                        Loading profile...
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="match-content">
            <div className="card-stack">
                {currentUser && (
                    <ProfileCard 
                        key={currentUser.id}
                        user={{
                            ...currentUser,
                            distance: distance ?? undefined
                        }}
                        onLike={handleLike}
                        onReject={handleReject}
                    />
                )}
            </div>
        </div>
    );
};

export default Home;
