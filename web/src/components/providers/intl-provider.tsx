'use client';

import { NextIntlClientProvider } from 'next-intl';
import { ReactNode, useEffect, useState } from 'react';
import { useLocale } from '@/hooks/use-locale';
import { defaultLocale } from '@/i18n';

interface IntlProviderProps {
  children: ReactNode;
}

export function IntlProvider({ children }: IntlProviderProps) {
  const { locale, isLoading: localeLoading } = useLocale();
  const [messages, setMessages] = useState<Record<string, unknown> | null>(null);
  const [messagesLoading, setMessagesLoading] = useState(true);

  useEffect(() => {
    if (localeLoading) return;

    async function loadMessages() {
      setMessagesLoading(true);
      try {
        const messageModule = await import(`@/messages/${locale}.json`);
        setMessages(messageModule.default);
      } catch (error) {
        console.error('Failed to load messages:', error);
        // use default locale as fallback
        const fallbackModule = await import(`@/messages/${defaultLocale}.json`);
        setMessages(fallbackModule.default);
      } finally {
        setMessagesLoading(false);
      }
    }

    loadMessages();
  }, [locale, localeLoading]);

  if (localeLoading || messagesLoading || !messages) {
    return <div className="flex items-center justify-center min-h-screen">Loading...</div>;
  }

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
}