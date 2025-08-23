import { Routes, Route, Navigate } from 'react-router-dom'
import { memo, type ReactNode } from 'react'
import { useAuthStore } from './features/auth/hooks/authStore'
import LoginForm from './features/auth/components/LoginForm'
import RegisterForm from './features/auth/components/RegisterForm'
import BottomNav from './shared/components/BottomNav'
import Header from './shared/components/Header'
import ViewProfile from './features/userProfile/components/ViewProfile'
import EditProfile from './features/userProfile/components/EditProfile'
import { useCurrentUser } from './features/userProfile/hooks/useCurrentUser'
import Home from './features/matches/components/Home'
import { useUserRecommendations } from './features/matches/hooks/useMatch'


const ProtectedRoute = memo(({ children }: { children: ReactNode }) => {
  const { authToken } = useAuthStore();
    const { isLoading } = useCurrentUser();
    useUserRecommendations()

  if (!authToken) {
    return <Navigate to="/login" replace />;
  }

  if (isLoading) {
    return <div>Loading...</div>; // Or a spinner component
  }


  return (
    <>
      <Header />
      {children}
      <BottomNav />
    </>
  );
});

ProtectedRoute.displayName = 'ProtectedRoute';


const AppRoutes = () => {

  return (
    <Routes>
      <Route path="/login" element={<LoginForm />} />
      <Route path="/register" element={<RegisterForm />} />

      {/* Protected Routes */}
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Home />
          </ProtectedRoute>
        }
      />
      <Route
        path="/home"
        element={
          <ProtectedRoute>
            <Home />
          </ProtectedRoute>
        }
      />
      <Route
        path="/profile"
        element={
          <ProtectedRoute>
            <ViewProfile />
          </ProtectedRoute>
        }
      />

      <Route
        path="/edit-profile"
        element={
          <ProtectedRoute>
            <EditProfile />
          </ProtectedRoute>
        }
      />


    </Routes>
  )
}

export default AppRoutes