import { memo, useState, useEffect } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { useAuthStore } from './features/auth/hooks/authStore'
import { useCurrentUser } from './features/userProfile/hooks/useCurrentUser'
import { useConnections, useConnectionRequests, } from './features/connections/hooks/useConnections'
import { useUserRecommendations } from './features/matches/hooks/useMatch'
import { useConnectionWebSocketEvents } from './features/connections/hooks/useConnectionWebSocketEvents'
import { useUIStore } from './shared/hooks/uiStore'
import Header from './shared/components/Header'
import BottomNav from './shared/components/BottomNav'
import IsLoading from './shared/components/IsLoading'

const ProtectedRoute = memo(() => {
  const { authToken } = useAuthStore()
  const { isLoading } = useCurrentUser()
  const { isChatMessageViewActive, view } = useUIStore()
  const [isMobile, setIsMobile] = useState(false)

  // Mobile detection
  useEffect(() => {
    const checkIsMobile = () => {
      setIsMobile(window.innerWidth <= 768)
    }

    checkIsMobile()
    window.addEventListener('resize', checkIsMobile)

    return () => window.removeEventListener('resize', checkIsMobile)
  }, [])

  // Run side-effect hooks
  useConnections()
  useUserRecommendations()
  useConnectionWebSocketEvents()
  useConnectionRequests()

  if (!authToken) {
    return <Navigate to="/login" replace />
  }

  if (isLoading) {
    return <IsLoading />
  }

  // Hide bottom nav on mobile when in chat message view
  const showBottomNav = !(isMobile && isChatMessageViewActive)
  // Hide header completely when in chat view
  const showHeader = view !== 'chat'

  return (
    <>
      {showHeader && <Header />}
      <Outlet />   {/* child route gets rendered here */}
      {showBottomNav && <BottomNav />}
    </>
  )
})

ProtectedRoute.displayName = 'ProtectedRoute'
export default ProtectedRoute
