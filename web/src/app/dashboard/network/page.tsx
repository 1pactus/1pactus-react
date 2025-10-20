'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useTranslations } from 'next-intl';
import { ChartLineInteractive } from './chart-line-interactive'
import { ChartLineLinear, ChartData } from './chart-line-linear';
import { GetNetworkHealthRequest, GetNetworkHealthResponse, NetworkStatusData } from '@/lib/proto/api/blockchain'
import { AlertCircleIcon, CheckCircle2Icon, PopcornIcon } from "lucide-react"
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/components/ui/alert"
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';

import React from 'react';

const API_MODE = process.env.NEXT_PUBLIC_API_MODE;
const ONE_PAC = 1000000000

const getApiBaseUrl = () => {
  const envUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '';
  if (envUrl.startsWith('/')) {
    return typeof window !== 'undefined' ? `${window.location.origin}${envUrl}` : envUrl;
  }
  return envUrl;
};

export type NetworkStatutsChartData = Array<{ 
    date: string; 

    stake: number;
    supply: number;
    circulating_supply: number;
    txs: number;
    blocks: number;
    fee: number;
    active_validator: number;
    active_account: number;
  }>;

function ChartLineContents({ chartData }: { chartData: NetworkStatutsChartData }) {
  const tNetworkOverview = useTranslations('network/overview');

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
      <div className='grid grid-cols-1 lg:grid-cols-2 gap-4'>
        {/*<ChartLineInteractive/>*/}
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.blocks }))}
          title={tNetworkOverview('block-committed-title')}
          description={tNetworkOverview('block-committed-description')}
          displayLable="Blocks"
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.txs }))}
          title={tNetworkOverview('transactions-committed-title')}
          description={tNetworkOverview('transactions-committed-description')}
          displayLable='txs'
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.stake }))}
          title={tNetworkOverview('stake-title')}
          description={tNetworkOverview('stake-description')}
          displayLable='stake'
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.supply }))}
          title={tNetworkOverview('supply-title')}
          description={tNetworkOverview('supply-description')}
          displayLable='supply'
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.circulating_supply }))}
          title={tNetworkOverview('circulating-supply-title')}
          description={tNetworkOverview('circulating-supply-description')}
          displayLable='circulating_supply'
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.fee }))}
          title={tNetworkOverview('fee-accumulation-title')}
          description={tNetworkOverview('fee-accumulation-description')}
          displayLable='fee'
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.active_validator }))}
          title={tNetworkOverview('active-validators-title')}
          description={tNetworkOverview('active-validators-description')}
          displayLable='active_validator'
        />
        <ChartLineLinear
          data={chartData.map(item => ({ date: item.date, value: item.active_account }))}
          title={tNetworkOverview('active-accounts-title')}
          description={tNetworkOverview('active-accounts-description')}
          displayLable='active_account'
        />
      </div>
    </div>
  )
}

export default function NetworkOverviewPage() {
  const [chartData, setChartData] = React.useState<NetworkStatutsChartData>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<boolean>(false);

  const tCommon = useTranslations('common');

  const fetchData = async () => {
      const reqPayload = GetNetworkHealthRequest.create({ days: -1, datatype: API_MODE === 'pb' ? 'pb' : 'json' });
      try {
        setIsLoading(true);
        setError(false);

        const apiBaseUrl = getApiBaseUrl();
        const jsonPayload = GetNetworkHealthRequest.toJSON(reqPayload) as Record<string, unknown>;
        const params = new URLSearchParams();

        Object.entries(jsonPayload).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== 0) {
            params.append(key, String(value));
          }
        });

        const response = await fetch(`${apiBaseUrl}/network_status?${params}`, {
          method: 'GET',
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        let ret: GetNetworkHealthResponse | null = null;

        switch (reqPayload.datatype) {
          case 'pb':
            // protocol buffer decode
            ret = GetNetworkHealthResponse.decode(new Uint8Array(await response.arrayBuffer()));
            break;
          case 'json':
            // JSON decode
            const jsonData = await response.json();
            ret = GetNetworkHealthResponse.fromJSON(jsonData);
            break;
          default:
            throw new Error(`Unsupported datatype: ${reqPayload.datatype}`);
        }

        if (ret.lines && ret.lines.length > 0) {
          const transformedData = ret.lines.map((item: NetworkStatusData) => ({
            date: new Date(item.timeIndex * 1000).toISOString().split('T')[0],
            stake: Number(item.stake.toString()) / ONE_PAC,
            supply: Number(item.supply.toString()) / ONE_PAC,
            circulating_supply: Number(item.circulatingSupply.toString()) / ONE_PAC,
            txs: Number(item.txs.toString()),
            blocks: Number(item.blocks.toString()),
            fee: Number(item.fee.toString()) / ONE_PAC,
            active_validator: Number(item.activeValidator.toString()),
            active_account: Number(item.activeAccount.toString()),
          }));
          setChartData(transformedData);
        }

        return ret;
      } catch (error) {
        console.error('Failed to fetch network health data:', error);
        setError(true);
      } finally {
        setIsLoading(false);
      }
    };

  React.useEffect(() => {
    fetchData();
  }, []);

  return (
      <>
      {isLoading ? (
        <div className="flex flex-1 items-center justify-center p-4">
          <div className="w-full max-w-md space-y-4">
            <p className="text-center text-muted-foreground">{tCommon('loading')}</p>
            <Progress value={undefined} className="w-full" />
          </div>
        </div>
      ) : error ? (
        <div className="flex flex-1 items-center justify-center p-4">
          <Alert variant="destructive" className="max-w-2xl">
          <AlertCircleIcon />
          <AlertTitle>{tCommon('failed-load')}</AlertTitle>
          <AlertDescription>
            <p>{tCommon('try-later')}</p>
            <Button onClick={() => fetchData()} >{tCommon('retry')}</Button>
          </AlertDescription>
        </Alert>
        </div>
      ) : (
        <>
          <ChartLineContents chartData={chartData} />
        </>
      )}
    </>
  );
}