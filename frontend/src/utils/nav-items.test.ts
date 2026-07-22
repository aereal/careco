import { describe, expect, it } from 'vitest';
import { isNavItemActive, navItems } from './nav-items';

describe('navItems', () => {
  it('contains a link to home', () => {
    expect(navItems).toContainEqual(
      expect.objectContaining({ type: 'link', label: 'ホーム', to: '/' }),
    );
  });

  it('contains an action item for picking a report month', () => {
    expect(navItems).toContainEqual(
      expect.objectContaining({
        type: 'action',
        label: '月別レポート',
        matchPrefix: '/reports',
      }),
    );
  });
});

describe('isNavItemActive', () => {
  const home = {
    type: 'link',
    label: 'ホーム',
    to: '/',
    matchPrefix: '/',
  } as const;
  const reports = {
    type: 'action',
    label: '月別レポート',
    matchPrefix: '/reports',
  } as const;

  it('marks home active only on the exact root path', () => {
    expect(isNavItemActive('/', home)).toBe(true);
    expect(isNavItemActive('/reports/2026/06', home)).toBe(false);
  });

  it('marks reports active for any month under /reports', () => {
    expect(isNavItemActive('/reports/2026/06', reports)).toBe(true);
    expect(isNavItemActive('/reports/2024/01', reports)).toBe(true);
    expect(isNavItemActive('/', reports)).toBe(false);
  });
});
