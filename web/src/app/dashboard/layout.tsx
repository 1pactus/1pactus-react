'use client'

import { AppSidebar } from "@/components/app-sidebar"
import { useNavData, NavPageRoot } from "@/hooks/use-nav-data"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"

import { Separator } from "@/components/ui/separator"

import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"

import { usePathname } from "next/navigation"
import Link from "next/link"
import { ThemeToggle } from "@/components/theme-toggle"
import { LanguageSwitcher } from "@/components/language-switcher"
import { useTranslations } from 'next-intl'

export default function Page({ children }: Readonly<{ children: React.ReactNode }>) {
  const pathname = usePathname()
  const t = useTranslations('navigation')
  const navData = useNavData()

  // Function to generate breadcrumb data
  const generateBreadcrumbs = () => {
    const breadcrumbs = [
      {
        title: t('dashboard'),
        href: NavPageRoot,
        isActive: false
      }
    ]

    // Match navigation items based on current path
    for (const navItem of navData.navMain) {
      for (const subItem of navItem.items || []) {
        if (subItem.url === pathname) {
          breadcrumbs.push({
            title: navItem.title,
            href: navItem.url,
            isActive: false
          })
          breadcrumbs.push({
            title: subItem.title,
            href: subItem.url,
            isActive: true
          })
          return breadcrumbs
        }
      }
    }

    // If no match is found, use default path parsing
    const segments = pathname.split('/').filter(Boolean)
    if (segments.length > 1) {
      // Remove 'i' prefix
      const pageSegments = segments.slice(1)
      pageSegments.forEach((segment, index) => {
        const isLast = index === pageSegments.length - 1
        breadcrumbs.push({
          title: segment.charAt(0).toUpperCase() + segment.slice(1),
          href: isLast ? pathname : `/${segments.slice(0, index + 2).join('/')}`,
          isActive: isLast
        })
      })
    }

    return breadcrumbs
  }

  const breadcrumbs = generateBreadcrumbs()
  return (
    <SidebarProvider>
      <AppSidebar data={navData} currentPath={pathname} />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 justify-between">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mr-2 data-[orientation=vertical]:h-4"
            />
            
            <Breadcrumb>
              <BreadcrumbList>
                {breadcrumbs.map((crumb, index) => (
                  <div key={crumb.href} className="flex items-center">
                    {index > 0 && (
                      <BreadcrumbSeparator className="mx-2" />
                    )}
                    <BreadcrumbItem>
                      {crumb.isActive ? (
                        <BreadcrumbPage>{crumb.title}</BreadcrumbPage>
                      ) : (
                        <BreadcrumbLink asChild>
                          <Link href={crumb.href}>{crumb.title}</Link>
                        </BreadcrumbLink>
                      )}
                    </BreadcrumbItem>
                  </div>
                ))}
              </BreadcrumbList>
            </Breadcrumb>
            
          </div>
          
          <div className="flex items-center gap-2 px-4">
            <LanguageSwitcher />
            <ThemeToggle />
          </div>
        </header>
        {children}
      </SidebarInset>
    </SidebarProvider>
  )
}
