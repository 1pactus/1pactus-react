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
      {/*
      <div className="grid auto-rows-min gap-4 md:grid-cols-3">
        <div className="bg-muted/50 aspect-video rounded-xl" />
        <div className="bg-muted/50 aspect-video rounded-xl" />
        <Card className="aspect-video">
          <CardHeader>
            <CardTitle>Internationalization Demo</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-4">
              Try out the multilingual support with English and Chinese. Notice how the sidebar navigation also switches languages!
            </p>
            <Link href="/intl-demo">
              <Button>View Demo</Button>
            </Link>
          </CardContent>
        </Card>
      </div>*/}
      <Card className="min-h-[100vh] flex-1 md:min-h-min">
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground mb-4">
            {t("description")}
          </p>
        </CardContent>

        <CardHeader>
          <CardTitle>{t("disclaimer-title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground mb-4">
            {t("disclaimer-content")}
          </p>
        </CardContent>
        <CardHeader>
          <CardTitle>{t("updates-title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <ul className="list-disc list-inside text-muted-foreground space-y-1 mt-2">
            <li>{t("updates-content-0")}</li>
          </ul>
        </CardContent>
      </Card>

      {/*<Card>
  <CardHeader>
    <CardTitle>Card Title</CardTitle>
    <CardDescription>Card Description</CardDescription>
    <CardAction>Card Action</CardAction>
  </CardHeader>
  <CardContent>
    <p>Card Content</p>
  </CardContent>
  <CardFooter>
    <p>Card Footer</p>
  </CardFooter>
</Card>*/}
    </div>
  )
}
