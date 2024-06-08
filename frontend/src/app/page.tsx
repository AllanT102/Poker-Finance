'use client'
import { useRouter, usePathname } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  const pathname = usePathname()
  const router = useRouter()

  useEffect(() => {
    if (pathname === '/') {
      router.push('/dashboard')
    }
  }, [pathname, router])

  return null;
}
