import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import type { UserRole } from '@/types'
import Navbar from './Navbar'

interface PrivateRouteProps {
  requiredRole?: UserRole
}

export default function PrivateRoute({ requiredRole }: PrivateRouteProps) {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requiredRole && user?.role !== requiredRole) {
    // Redirect to the user's own dashboard
    const fallback = user?.role === 'recruiter' ? '/recruiter/jobs' : '/candidate/jobs'
    return <Navigate to={fallback} replace />
  }

  return (
    <>
      <Navbar />
      <main className="container mx-auto px-4 py-8 max-w-6xl">
        <Outlet />
      </main>
    </>
  )
}
