import { Link, useNavigate } from 'react-router-dom'
import { LogOut, Briefcase, ClipboardList } from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/button'

export default function Navbar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <header className="border-b bg-white shadow-sm">
      <div className="container mx-auto px-4 max-w-6xl flex items-center justify-between h-14">
        <Link to="/" className="font-semibold text-primary text-lg tracking-tight">
          R&amp;S
        </Link>

        <nav className="flex items-center gap-4">
          {user?.role === 'recruiter' && (
            <>
              <Link
                to="/recruiter/jobs"
                className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <Briefcase className="w-4 h-4" />
                Minhas Vagas
              </Link>
            </>
          )}

          {user?.role === 'candidate' && (
            <>
              <Link
                to="/candidate/jobs"
                className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <Briefcase className="w-4 h-4" />
                Vagas
              </Link>
              <Link
                to="/candidate/applications"
                className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <ClipboardList className="w-4 h-4" />
                Minhas Candidaturas
              </Link>
            </>
          )}

          <span className="text-sm text-muted-foreground hidden sm:block">{user?.name || user?.email}</span>

          <Button variant="ghost" size="sm" onClick={handleLogout} className="gap-1.5">
            <LogOut className="w-4 h-4" />
            <span className="hidden sm:inline">Sair</span>
          </Button>
        </nav>
      </div>
    </header>
  )
}
