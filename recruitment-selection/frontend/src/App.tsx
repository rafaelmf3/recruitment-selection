import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import PrivateRoute from '@/components/PrivateRoute'
import PublicRoute from '@/components/PublicRoute'
import LoginPage from '@/pages/auth/LoginPage'
import RegisterPage from '@/pages/auth/RegisterPage'
import MyJobsPage from '@/pages/recruiter/MyJobsPage'
import CreateJobPage from '@/pages/recruiter/CreateJobPage'
import JobPipelinePage from '@/pages/recruiter/JobPipelinePage'
import BrowseJobsPage from '@/pages/candidate/BrowseJobsPage'
import MyApplicationsPage from '@/pages/candidate/MyApplicationsPage'

function DashboardRedirect() {
  const { user } = useAuth()
  if (user?.role === 'recruiter') return <Navigate to="/recruiter/jobs" replace />
  return <Navigate to="/candidate/jobs" replace />
}

export default function App() {
  return (
    <Routes>
      {/* Public-only routes (redirect to dashboard if already authenticated) */}
      <Route element={<PublicRoute />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Route>

      {/* Recruiter routes */}
      <Route element={<PrivateRoute requiredRole="recruiter" />}>
        <Route path="/recruiter/jobs" element={<MyJobsPage />} />
        <Route path="/recruiter/jobs/new" element={<CreateJobPage />} />
        <Route path="/recruiter/jobs/:id/pipeline" element={<JobPipelinePage />} />
      </Route>

      {/* Candidate routes */}
      <Route element={<PrivateRoute requiredRole="candidate" />}>
        <Route path="/candidate/jobs" element={<BrowseJobsPage />} />
        <Route path="/candidate/applications" element={<MyApplicationsPage />} />
      </Route>

      {/* Default redirect */}
      <Route path="/" element={<PrivateRoute />}>
        <Route index element={<DashboardRedirect />} />
      </Route>

      {/* Catch-all */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
