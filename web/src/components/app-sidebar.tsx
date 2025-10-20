"use client"

import * as React from "react"
import {
  ChartNoAxesCombined
} from "lucide-react"

import Link from "next/link"
import { NavMain, NavMainProps } from "@/components/nav-main"
import { NavProjects, NavProjectsProps } from "@/components/nav-projects"
import { NavSecondary, NavSecondaryProps } from "@/components/nav-secondary"
import { NavUser, NavUserProps } from "@/components/nav-user"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

export interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  data: {
    homeUrl?: string;
    user?: NavUserProps["user"];
    navMain: NavMainProps["items"];
    navSecondary: NavSecondaryProps["items"];
    projects: NavProjectsProps["projects"];
  };
  currentPath?: string;
}

export function AppSidebar({ data, currentPath, ...props }: AppSidebarProps) {
  return (
    <Sidebar variant="inset" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link href={ data.homeUrl || "/" }>
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <ChartNoAxesCombined className="size-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">1PACTUS</span>
                  <span className="truncate text-xs">Charts</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} currentPath={currentPath} />
        {/*<NavProjects projects={data.projects} />*/}
        {/*<NavSecondary items={data.navSecondary} className="mt-auto" />*/}
      </SidebarContent>
      {/*<SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>*/}
    </Sidebar>
  )
}
