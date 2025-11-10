import { atom } from 'jotai';

export const distance = atom(0);
export const date = atom(new Date());
export const active = atom(true);
export const memo = atom<string>();
