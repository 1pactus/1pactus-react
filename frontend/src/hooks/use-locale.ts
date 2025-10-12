'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Locale, defaultLocale, locales } from '@/i18n';

const LOCALE_STORAGE_KEY = 'preferred-locale';

interface LocaleContextType {
  locale: Locale;
  changeLocale: (newLocale: Locale) => void;
  isLoading: boolean;
}

const LocaleContext = createContext<LocaleContextType | undefined>(undefined);

export function LocaleProvider({ children }: { children: ReactNode }) {
  const [locale, setLocale] = useState<Locale>(defaultLocale);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const savedLocale = localStorage.getItem(LOCALE_STORAGE_KEY) as Locale;
    if (savedLocale && locales.includes(savedLocale)) {
      setLocale(savedLocale);
    }
    setIsLoading(false);
  }, []);

  const changeLocale = (newLocale: Locale) => {
    setLocale(newLocale);
    localStorage.setItem(LOCALE_STORAGE_KEY, newLocale);
  };

  const value = React.useMemo(() => ({
    locale,
    changeLocale,
    isLoading
  }), [locale, isLoading]);

  return React.createElement(
    LocaleContext.Provider,
    { value },
    children
  );
}

export function useLocale() {
  const context = useContext(LocaleContext);
  if (context === undefined) {
    throw new Error('useLocale must be used within a LocaleProvider');
  }
  return context;
}