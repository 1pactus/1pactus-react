'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useTranslations } from 'next-intl';
import { ChartLineInteractive } from './chart-line-interactive'
import { ChartLineLinear } from './chart-line-linear';
import { GetNetworkHealthRequest, GetNetworkHealthResponse } from '@/lib/proto/api/blockchain'
import React from 'react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;
const API_MODE = process.env.NEXT_PUBLIC_API_MODE;

export default function NetworkOverviewPage() {
  const t = useTranslations('navigation');

  React.useEffect(() => {
    const reqPayload = GetNetworkHealthRequest.create({ days: 365, datatype: API_MODE === 'pb' ? 'pb' : 'json' });

    const fetchData = async () => {
      const jsonPayload = GetNetworkHealthRequest.toJSON(reqPayload) as any;
      const params = new URLSearchParams();

      Object.entries(jsonPayload).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== 0) {
          params.append(key, String(value));
        }
      });

      const response = await fetch(`${API_BASE_URL}/network_status?${params}`, {
        method: 'GET',
      });

      switch (reqPayload.datatype) {
        case 'pb':
          // TODO
        case 'json':
          const data: GetNetworkHealthResponse = await response.json();
        console.log(data);
      }

      
    };

    fetchData();
  }, []);

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
      {/*
      <Card>
        <CardHeader>
          <CardTitle>{t('network')} - {t('overview')}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            This is a demo page for the Network Overview section. 
            Notice how the breadcrumb navigation and sidebar are fully localized!
          </p>
          <p className="text-muted-foreground mt-4">
            <strong>Current Features:</strong>
          </p>
          <ul className="list-disc list-inside text-muted-foreground space-y-1 mt-2">
            <li>Multilingual sidebar navigation</li>
            <li>Localized breadcrumb navigation</li>
            <li>Automatic language detection</li>
            <li>Real-time language switching</li>
          </ul>
        </CardContent>
      </Card>*/}

      {/*<ChartLineInteractive/>*/}
      <ChartLineLinear/>
    </div>
  );
}