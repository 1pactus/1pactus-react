'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useTranslations } from 'next-intl';

export default function NetworkOverviewPage() {
  const t = useTranslations('navigation');

  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
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
      </Card>
    </div>
  );
}