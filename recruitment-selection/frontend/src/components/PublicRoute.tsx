import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'

export default function PublicRoute() {
  const { isAuthenticated, user } = useAuth()

  if (isAuthenticated) {
    const destination = user?.role === 'recruiter' ? '/recruiter/jobs' : '/candidate/jobs'
    return <Navigate to={destination} replace />
  }

  return <Outlet />
}
