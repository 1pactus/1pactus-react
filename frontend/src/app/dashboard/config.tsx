import {
  LifeBuoy,
  Map,
  Send,
  UserRoundSearch,
  Trophy,
  TrendingUp,
} from "lucide-react"

export const NavPageRoot = "/dashboard"

export const NavData = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Network",
      url: `#`,
      icon: Map,
      items: [
        {
          title: "Overview",
          url: `${NavPageRoot}/network/overview`,
        },
        {
          title: "Validators",
          url: `${NavPageRoot}/network/validators`,
        },
      ],
    },
    {
      title: "Trends",
      url: `#`,
      icon: TrendingUp,
      items: [
        {
          title: "Supply",
          url: `${NavPageRoot}/trends/supply`,
        },
        {
          title: "Bond",
          url: `${NavPageRoot}/trends/bond`,
        },
        {
          title: "Unbond",
          url: `${NavPageRoot}/trends/unbond`,
        },
        {
          title: "Stake",
          url: `${NavPageRoot}/trends/stake`,
        },
        {
          title: "Reward",
          url: `${NavPageRoot}/trends/reward`,
        }
      ],
    },
    {
      title: "Top Addresses",
      url: `#`,
      icon: Trophy,
      items: [
        {
          title: "Top Transfer",
          url: `${NavPageRoot}/top_addresses/transfer`,
        },
        {
          title: "Top Bond",
          url: `${NavPageRoot}/top_addresses/bond`,
        },
        {
          title: "Top Reward",
          url: `${NavPageRoot}/top_addresses/reward`,
        },
        {
          title: "Top Withdraw",
          url: `${NavPageRoot}/top_addresses/withdraw`,
        }
      ],
    },
    {
      title: "Public Addresses",
      url: `${NavPageRoot}/public_addresses`,
      icon: UserRoundSearch,
    }
  ],
  navSecondary: [
    {
      title: "Support",
      url: "#",
      icon: LifeBuoy,
    },
    {
      title: "Feedback",
      url: "#",
      icon: Send,
    },
  ],
  projects: [
    {
      name: "Watch Address",
      url: `${NavPageRoot}/watch_address`,
      icon: UserRoundSearch,
    },
  ],
}

export const CoinUnit = "$PAC"