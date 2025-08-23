import { Routes, Route, Navigate } from 'react-router-dom'
import TestRoute from './features/userProfile/components/Test'
import { memo, type ReactNode } from 'react'
import { useAuthStore } from './features/auth/hooks/authStore'
import LoginForm from './features/auth/components/LoginForm'
import RegisterForm from './features/auth/components/RegisterForm'
import BottomNav from './shared/components/BottomNav'
import Header from './shared/components/Header'
import EditProfile from './features/userProfile/components/EditProfile'
import { useCurrentUser } from './features/userProfile/hooks/useCurrentUser'

const ProtectedRoute = memo(({ children }: { children: ReactNode }) => {
  const { authToken } = useAuthStore();
    const { isLoading } = useCurrentUser();

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
      {/* <Route path="/" element={<LoginForm />} /> */}
      <Route path="/login" element={<LoginForm />} />
      <Route path="/register" element={<RegisterForm />} />

      {/* Protected Routes */}
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <TestRoute />
          </ProtectedRoute>
        }
      />

      <Route
        path="/profile"
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