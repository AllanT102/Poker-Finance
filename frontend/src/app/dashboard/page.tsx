'use client'
import { useUser } from "@clerk/nextjs";
import MainLayout from "./layout"


const sections = [
  { name: 'Dashboard', href: '/dashboard' },
  { name: 'Profile', href: '/profile' },
  { name: 'Settings', href: '/settings' },
]

const Dashboard = () => {
  const { isLoaded, isSignedIn, user } = useUser();
  return (
    <div>
      <h1 className="text-2xl font-bold">Welcome to the Home Page</h1>
      <p>This is the home page content.</p>
      <div>Hello, {user?.firstName} welcome to Clerk</div>
    </div>
  )
}

export default Dashboard