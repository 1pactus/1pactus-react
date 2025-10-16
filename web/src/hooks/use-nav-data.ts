'use client';

import {
  LifeBuoy,
  Map,
  Send,
  UserRoundSearch,
  Trophy,
  TrendingUp,
} from "lucide-react";
import { useTranslations } from 'next-intl';

export const NavPageRoot = "/dashboard";

export function useNavData() {
  const t = useTranslations('navigation');

  return {
    navMain: [
      {
        title: t('network'),
        url: `#`,
        icon: Map,
        items: [
          {
            title: t('overview'),
            url: `${NavPageRoot}/network/overview`,
          },
          {
            title: t('validators'),
            url: `${NavPageRoot}/network/validators`,
          },
        ],
      },
      {
        title: t('trends'),
        url: `#`,
        icon: TrendingUp,
        items: [
          {
            title: t('supply'),
            url: `${NavPageRoot}/trends/supply`,
          },
          {
            title: t('bond'),
            url: `${NavPageRoot}/trends/bond`,
          },
          {
            title: t('unbond'),
            url: `${NavPageRoot}/trends/unbond`,
          },
          {
            title: t('stake'),
            url: `${NavPageRoot}/trends/stake`,
          },
          {
            title: t('reward'),
            url: `${NavPageRoot}/trends/reward`,
          }
        ],
      },
      {
        title: t('topAddresses'),
        url: `#`,
        icon: Trophy,
        items: [
          {
            title: t('topTransfer'),
            url: `${NavPageRoot}/top_addresses/transfer`,
          },
          {
            title: t('topBond'),
            url: `${NavPageRoot}/top_addresses/bond`,
          },
          {
            title: t('topReward'),
            url: `${NavPageRoot}/top_addresses/reward`,
          },
          {
            title: t('topWithdraw'),
            url: `${NavPageRoot}/top_addresses/withdraw`,
          }
        ],
      },
      {
        title: t('publicAddresses'),
        url: `${NavPageRoot}/public_addresses`,
        icon: UserRoundSearch,
      }
    ],
    navSecondary: [
      {
        title: t('support'),
        url: "#",
        icon: LifeBuoy,
      },
      {
        title: t('feedback'),
        url: "#",
        icon: Send,
      },
    ],
    projects: [
      {
        name: t('watchAddress'),
        url: `${NavPageRoot}/watch_address`,
        icon: UserRoundSearch,
      },
    ],
  };
}

export const CoinUnit = "$PAC";