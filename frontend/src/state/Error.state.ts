import { atom, selector } from "recoil";

export const errorState = atom<string | null>({
  key: "errorState",
  default: null,
});

export const errorStateSelector = selector({
  key: "getError",
  get: ({ get }) => {
    const error = get(errorState);
    return error;
  },
  set: ({set},newState) => {
    set(errorState, newState);
  }
});