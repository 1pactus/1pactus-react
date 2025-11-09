"use client"

import { TrendingUp } from "lucide-react"
import { CartesianGrid, Line, LineChart, XAxis, Bar, BarChart } from "recharts"
import { useTranslations } from 'next-intl';

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"

export const description = "A linear line chart"

export type ChartData = Array<{ 
    date: string; 

    value: number;
}>

const getChartConfig = (unitLabel: string): ChartConfig => {
  const defaultLabel = unitLabel || "value";
  return {
    ["value"]: {
      label: ` ${defaultLabel}`,
      color: "var(--chart-1)",
    },
  } satisfies ChartConfig;
}

export type ChartLineLinearProps = {
  data: ChartData;
  title: string;
  description: string;
  displayLable: string;
  footerChildren?: React.ReactNode;
  barThreshold?: number; // 使用柱状图的阈值，默认为60
}

export function ChartLineLinear({ 
  data, 
  title, 
  description, 
  displayLable: displayLable, 
  footerChildren,
  barThreshold = 30
}: ChartLineLinearProps) {
  const chartConfig = getChartConfig(displayLable);
  const useBarChart = data.length <= barThreshold;
  
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="aspect-auto h-[250px] w-full">
          {useBarChart ? (
            <BarChart
              accessibilityLayer
              data={data}
              margin={{
                left: 12,
                right: 12,
              }}
            >
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="date"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                minTickGap={32}
                tickFormatter={(value) => {
                  const date = new Date(value)
                  return date.toISOString().split('T')[0]
                }}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent 
                    hideLabel={false}
                    labelFormatter={(value) => {
                      const date = new Date(value);
                      return date.toISOString().split('T')[0];
                    }}
                    formatter={(value) => [
                      value?.toLocaleString(),
                      ""
                    ]}
                  />
                }
              />
              <Bar 
                dataKey="value" 
                fill={`var(--color-value)`}
                radius={[4, 4, 0, 0]}
                isAnimationActive={false}
              />
            </BarChart>
          ) : (
            <LineChart
              accessibilityLayer
              data={data}
              margin={{
                left: 12,
                right: 12,
              }}
            >
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="date"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                minTickGap={32}
                tickFormatter={(value) => {
                  const date = new Date(value)
                  return date.toISOString().split('T')[0]
                }}
              />
              <ChartTooltip
                cursor={false}
                content={
                  <ChartTooltipContent 
                    hideLabel={false}
                    labelFormatter={(value) => {
                      const date = new Date(value);
                      return date.toISOString().split('T')[0];
                    }}
                    formatter={(value) => [
                      value?.toLocaleString(),
                      ""
                    ]}
                  />
                }
              />
              <Line
                dataKey="value"
                type="linear"
                stroke={`var(--color-value)`}
                strokeWidth={2}
                dot={false}
                isAnimationActive={false}
              />
            </LineChart>
          )}
        </ChartContainer>
      </CardContent>
      {
        footerChildren ? (
          <CardFooter>
            {footerChildren}
          </CardFooter>
        ) : null
      }
      {/*<CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 leading-none font-medium">
          Trending up by 5.2% this month <TrendingUp className="h-4 w-4" />
        </div>
        <div className="text-muted-foreground leading-none">
          Showing total visitors for the last 6 months
        </div>
      </CardFooter>*/}
    </Card>
  )
}
