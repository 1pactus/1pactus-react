"use client"

import Link from "next/link"

import * as React from "react"
import { ChevronRight, type LucideIcon } from "lucide-react"

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"

export interface NavMainProps {
  items: {
    title: string
    url: string
    icon: LucideIcon
    isActive?: boolean
    items?: {
      title: string
      url: string
    }[]
  }[]
  currentPath?: string
}

export function NavMain({ items, currentPath }: NavMainProps) {
  const [openState, setOpenState] = React.useState<Record<string, boolean>>({})

  React.useEffect(() => {
    const initialState: Record<string, boolean> = {}
    for (const item of items) {
      const isParentActive = currentPath === item.url
      const hasActiveSub =
        item.items?.some((subItem) => subItem.url === currentPath) || false
      if (isParentActive || hasActiveSub) {
        initialState[item.title] = true
      }
    }
    setOpenState(initialState)
  }, [currentPath, items])

  // 检查子菜单项是否为当前活跃页面
  const isSubItemActive = (subItemUrl: string) => {
    return currentPath === subItemUrl
  }

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Platform</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => {
          const isParentActive = currentPath === item.url
          return (
            <Collapsible
              key={item.title}
              asChild
              open={openState[item.title] || false}
              onOpenChange={(isOpen) =>
                setOpenState((prevState) => ({
                  ...prevState,
                  [item.title]: isOpen,
                }))
              }
            >
              <SidebarMenuItem>
                <SidebarMenuButton
                  asChild
                  tooltip={item.title}
                  isActive={isParentActive}
                >
                  <Link href={item.url}>
                    <item.icon />
                    <span>{item.title}</span>
                  </Link>
                </SidebarMenuButton>
                {item.items?.length ? (
                  <>
                    <CollapsibleTrigger asChild>
                      <SidebarMenuAction className="data-[state=open]:rotate-90">
                        <ChevronRight />
                        <span className="sr-only">Toggle</span>
                      </SidebarMenuAction>
                    </CollapsibleTrigger>
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        {item.items?.map((subItem) => (
                          <SidebarMenuSubItem key={subItem.title}>
                            <SidebarMenuSubButton
                              asChild
                              isActive={isSubItemActive(subItem.url)}
                            >
                              <Link href={subItem.url}>
                                <span>{subItem.title}</span>
                              </Link>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        ))}
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  </>
                ) : null}
              </SidebarMenuItem>
            </Collapsible>
          )
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}
