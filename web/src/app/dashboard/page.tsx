import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function Page() {
  return (
    <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
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
      </div>
      <Card className="min-h-[100vh] flex-1 md:min-h-min">
        <CardHeader>
          <CardTitle>Multi-language Sidebar Demo</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground mb-4">
            The sidebar navigation is now fully internationalized! Try switching between English and Chinese using the language switcher in the top-right corner.
          </p>
          <p className="text-muted-foreground">
            <strong>Features:</strong>
          </p>
          <ul className="list-disc list-inside text-muted-foreground space-y-1 mt-2">
            <li>Sidebar navigation titles translate instantly</li>
            <li>Breadcrumb navigation is also localized</li>
            <li>Language preference is saved in localStorage</li>
            <li>No page refresh required for language switching</li>
          </ul>
        </CardContent>
      </Card>
    </div>
  )
}
