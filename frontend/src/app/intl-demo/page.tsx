'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { LanguageSwitcher } from '@/components/language-switcher';
import { useTranslations } from 'next-intl';
import { useLocale } from '@/hooks/use-locale';
import { languageLabels } from '@/i18n';
import { useState } from 'react';

export default function IntlDemoPage() {
  const t = useTranslations('demo');
  const tCommon = useTranslations('common');
  const { locale, changeLocale } = useLocale();
  const [clickCount, setClickCount] = useState(0);

  return (
    <div className="container mx-auto p-6 space-y-6">
      {/* 顶部导航栏 */}
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">{t('title')}</h1>
        <LanguageSwitcher />
      </div>

      {/* 当前语言状态 */}
      <Card>
        <CardHeader>
          <CardTitle>{t('currentLanguage')}</CardTitle>
          <CardDescription>{t('description')}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-2">
            <span>{tCommon('language')}:</span>
            <Badge variant="secondary">{languageLabels[locale]}</Badge>
          </div>
        </CardContent>
      </Card>

      {/* 快速切换测试 */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Language Switch Test / 快速语言切换测试</CardTitle>
          <CardDescription>
            Click the buttons below to quickly test language switching / 点击下方按钮快速测试语言切换
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2">
            <Button 
              onClick={() => changeLocale('en')} 
              variant={locale === 'en' ? 'default' : 'outline'}
            >
              Switch to English
            </Button>
            <Button 
              onClick={() => changeLocale('zh')} 
              variant={locale === 'zh' ? 'default' : 'outline'}
            >
              切换到中文
            </Button>
          </div>
          <div className="p-4 bg-muted rounded-lg">
            <p className="font-semibold">{t('content.greeting')}</p>
            <p className="text-sm text-muted-foreground mt-2">
              Current locale: {locale} | Click count: {clickCount}
            </p>
            <Button 
              onClick={() => setClickCount(c => c + 1)} 
              className="mt-2"
              size="sm"
            >
              Test Click Counter
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* 示例内容 */}
      <Card>
        <CardHeader>
          <CardTitle>{tCommon('welcome')}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-lg">{t('content.greeting')}</p>
          
          <div>
            <h3 className="text-lg font-semibold mb-2">{t('content.features')}</h3>
            <ul className="list-disc list-inside space-y-1">
              <li>{t('content.feature1')}</li>
              <li>{t('content.feature2')}</li>
              <li>{t('content.feature3')}</li>
            </ul>
          </div>
        </CardContent>
      </Card>

      {/* 使用说明 */}
      <Card>
        <CardHeader>
          <CardTitle>Usage Instructions / 使用说明</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <p>
            <strong>English:</strong> Click the language switcher in the top-right corner to change the interface language.
          </p>
          <p>
            <strong>中文:</strong> 点击右上角的语言切换器来更改界面语言。
          </p>
          <p className="text-sm text-muted-foreground mt-4">
            <strong>Real-time test:</strong> The content should update immediately when you change the language, without requiring a page refresh.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}