// 支持的语言列表
export const locales = ['en', 'zh'] as const;
export type Locale = typeof locales[number];

// 默认语言
export const defaultLocale: Locale = 'en';

// 语言标签映射
export const languageLabels = {
  en: 'English',
  zh: '中文'
} as const;