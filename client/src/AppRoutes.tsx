import { Routes, Route } from 'react-router-dom'
import LoginForm from './features/auth/components/LoginForm'
import RegisterForm from './features/auth/components/RegisterForm'
import ViewProfile from './features/userProfile/components/ViewProfile'
import EditProfile from './features/userProfile/components/EditProfile'
import Home from './features/matches/components/Home'
import Connection from './features/connections/components/Connections'
import ProtectedRoute from './ProtectedRoute'
import ChatInterface from './features/chat/components/ChatInterface'
import ViewOtherProfile from './features/userProfile/components/ViewOtherProfile'
  

const AppRoutes = () => {
  return (
    <Routes>
      {/* Public Routes */}
      <Route path="/login" element={<LoginForm />} />
      <Route path="/register" element={<RegisterForm />} />

      {/* Protected Routes */}
      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<Home />} />
        <Route path="/home" element={<Home />} />
        <Route path="/users/:id" element={<ViewOtherProfile />} />
        <Route path="/profile" element={<ViewProfile />} />
        <Route path="/edit-profile" element={<EditProfile />} />
        <Route path="/connections" element={<Connection />} />
        <Route path="/chat" element={<ChatInterface />} />
      </Route>
    </Routes>
  )
}

export default AppRoutes
