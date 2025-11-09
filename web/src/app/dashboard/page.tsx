'use client';

import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { 
  Card, 
  CardContent, 
  CardHeader, 
  CardDescription,
  CardAction,
  CardTitle,
  CardFooter,
} from '@/components/ui/card';

import { useTranslations } from 'next-intl';

export default function Page() {
  const t = useTranslations('home');

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground mb-4">
            {t("description")}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t("disclaimer-title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground mb-4">
            {t("disclaimer-content")}
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t("updates-title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <ul className="list-disc list-inside text-muted-foreground space-y-1 mt-2">
            <li>{t("updates-content-0")}</li>
          </ul>
        </CardContent>
      </Card>
    </div>
  )
}
