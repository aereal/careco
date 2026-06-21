export type NavItem =
  | { type: 'link'; label: string; to: string; matchPrefix: string }
  | { type: 'action'; label: string; matchPrefix: string };

export const navItems: NavItem[] = [
  { type: 'link', label: 'ホーム', to: '/', matchPrefix: '/' },
  { type: 'action', label: '月別レポート', matchPrefix: '/reports' },
];

export const isNavItemActive = (currentPath: string, item: NavItem): boolean =>
  item.matchPrefix === '/'
    ? currentPath === '/'
    : currentPath.startsWith(item.matchPrefix);
