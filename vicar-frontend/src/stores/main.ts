import {defineStore} from "pinia";
import type {User} from "../@types/user.ts";
import {getBrowserName} from "@/utils/browser.ts";
import {getUser} from "@/rest/api/users.ts";

export const useMainStore = defineStore('main', {
  state: () => ({
    user: null as unknown as User,
    deviceName: null as string|null,
  }),
  actions: {
    async loadUser() {
      const user = await getUser();
      this.user = user as unknown as User;

      return user;
    },
    async getDeviceName(): Promise<string> {
      if (this.deviceName) {
        return this.deviceName;
      }

      try {
        const {ip} = await (await fetch("https://api.ipify.org/?format=json")).json();
        const {country_name, city} = await (await fetch(`https://ipapi.co/${ip}/json`)).json();
        const devcName = getBrowserName() + " (" + city + ", " + country_name + " - " + ip + ")";
        for (const [umlaut, ascii] of Object.entries({"ä": "ae", "ö": "oe", "ü": "ue", "Ä": "Ae", "Ö": "Oe", "Ü": "Ue", "ß": "ss"})) {
          devcName.replace(new RegExp(umlaut, "g"), ascii);
        }
        devcName.replace(/[^a-zA-Z0-9 .,_()-]/g, "");
        return this.deviceName = devcName;
      } catch (error) {
        console.error("Failed to fetch IP address");
        return getBrowserName();
      }
    },
  },
  getters: {
    hasTeam(): boolean {
      return this.user.isTeam;
    }
  }
});