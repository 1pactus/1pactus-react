'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useTranslations } from 'next-intl';
import { ChartLineInteractive } from './chart-line-interactive'
import { ChartLineLinear, ChartData } from './chart-line-linear';
import { GetNetworkHealthRequest, GetNetworkHealthResponse, NetworkStatusData } from '@/lib/proto/api/blockchain'
import React from 'react';

const API_MODE = process.env.NEXT_PUBLIC_API_MODE;
const ONE_PAC = 1000000000

// 处理 API_BASE_URL,如果是相对路径则拼接当前域名
const getApiBaseUrl = () => {
  const envUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '';
  if (envUrl.startsWith('/')) {
    // 如果是相对路径,拼接当前页面的 origin (协议 + 域名 + 端口)
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
  }>

export default function NetworkOverviewPage() {
  const [chartData, setChartData] = React.useState<NetworkStatutsChartData>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const reqPayload = GetNetworkHealthRequest.create({ days: -1, datatype: API_MODE === 'pb' ? 'pb' : 'json' });

    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

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
            // Protocol Buffer 解码
            ret = GetNetworkHealthResponse.decode(new Uint8Array(await response.arrayBuffer()));
            break;
          case 'json':
            // JSON 解码
            const jsonData = await response.json();
            ret = GetNetworkHealthResponse.fromJSON(jsonData);
            break;
          default:
            throw new Error(`Unsupported datatype: ${reqPayload.datatype}`);
        }

        console.log('Network health data:', ret);

        // 转换API数据为图表数据格式
        if (ret.lines && ret.lines.length > 0) {
          const transformedData = ret.lines.map((item: NetworkStatusData) => ({
            date: new Date(item.timeIndex * 1000).toISOString().split('T')[0], // 假设timeIndex是Unix时间戳
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
        setError(error instanceof Error ? error.message : '获取网络数据失败');
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  const tNetworkOverview = useTranslations('network/overview');

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

      {isLoading ? (
        <Card>
          <CardContent className="flex items-center justify-center h-64">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
              <p className="text-muted-foreground">Loading...</p>
            </div>
          </CardContent>
        </Card>
      ) : error ? (
        <Card>
          <CardContent className="flex items-center justify-center h-64">
            <div className="text-center">
              <p className="text-red-500 mb-2">错误</p>
              <p className="text-muted-foreground text-sm">{error}</p>
            </div>
          </CardContent>
        </Card>
      ) : (
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
      )}
    </div>
  );
}