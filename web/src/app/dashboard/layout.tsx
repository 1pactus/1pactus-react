'use client'

import { AppSidebar, AppSidebarProps } from "@/components/app-sidebar"
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
import { Button } from "@/components/ui/button"

export default function Page({ children }: Readonly<{ children: React.ReactNode }>) {
  const pathname = usePathname()
  const t = useTranslations('navigation')
  const navData = useNavData() as AppSidebarProps["data"];

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

        const title = segment.charAt(0).toLocaleLowerCase() + segment.slice(1)

        breadcrumbs.push({
          title: t(title),
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
            <Button
              variant="ghost"
              size="icon"
              asChild
            >
              <a
                href="https://github.com/1pactus/1pactus-react"
                target="_blank"
                rel="noopener noreferrer"
                aria-label="GitHub"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                  className="h-5 w-5"
                >
                  <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                </svg>
              </a>
            </Button>
          </div>
        </header>
        {children}
      </SidebarInset>
    </SidebarProvider>
  )
}
