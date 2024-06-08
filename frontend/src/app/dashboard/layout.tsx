'use client'
import { ReactNode, useState } from 'react'
import { Nav } from '../../components/ui/sidenav'
import {
  AlertCircle,
  Archive,
  ArchiveX,
  File,
  Inbox,
  MessagesSquare,
  PenBox,
  Search,
  Send,
  ShoppingCart,
  Trash2,
  Users2,
} from "lucide-react"

type MainLayoutProps = {
  children: ReactNode
  sections: { name: string; href: string }[]
}

const MainLayout: React.FC<MainLayoutProps> = ({ children, sections }) => {
  const [isCollapsed, setIsCollapsed] = useState(false)

  return (
    <div className="flex">
      <Nav 
        isCollapsed={isCollapsed}
        links={[
              {
                title: "Inbox",
                label: "128",
                icon: Inbox,
                variant: "default",
              },
              {
                title: "Drafts",
                label: "9",
                icon: File,
                variant: "ghost",
              },
              {
                title: "Sent",
                label: "",
                icon: Send,
                variant: "ghost",
              },
              {
                title: "Junk",
                label: "23",
                icon: ArchiveX,
                variant: "ghost",
              },
              {
                title: "Trash",
                label: "",
                icon: Trash2,
                variant: "ghost",
              },
              {
                title: "Archive",
                label: "",
                icon: Archive,
                variant: "ghost",
              },
            ]}/>
      <main className="flex-grow p-4">
        {children}
      </main>
    </div>
  )
}

export default MainLayout