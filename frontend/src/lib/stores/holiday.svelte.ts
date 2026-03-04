import { get, post } from "$lib/api";

interface Holiday {
  date: string;
  localName: string;
  name: string;
}

interface Lookup {
  id: string;
  holidayName: string;
  holidayDate: string;
  daysUntil: number;
  createdAt: string;
}

let holidays = $state<Holiday[]>([]);
let lookups = $state<Lookup[]>([]);

export function getHolidays() {
  return holidays;
}

export function getLookups() {
  return lookups;
}

export async function loadHolidays() {
  holidays = await get("/holidays");
}

export async function loadLookups() {
  lookups = await get("/lookups");
}

export async function calculateDays(holidayName: string, holidayDate: string) {
  const lookup = await post("/lookups", { holidayName, holidayDate });
  lookups = [lookup, ...lookups];
}
