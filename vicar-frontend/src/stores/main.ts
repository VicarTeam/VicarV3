import {defineStore} from "pinia";
import {Role, User} from "../@types/user.ts";
import {getBrowserName} from "@/utils/browser.ts";
import {getUser, getUserGroups} from "@/rest/api/users.ts";
import {FullGroup, Group} from "@/@types/group.ts";

export const useMainStore = defineStore('main', {
  state: () => ({
    user: null as unknown as User,
    groups: [] as Group[],
    userOptionsVisible: false as boolean,
    deviceName: null as string|null,
  }),
  actions: {
    async loadUser() {
      const user = await getUser();
      this.user = user as unknown as User;

      if (!user) {
        return false;
      }

      this.groups = await getUserGroups();
      return true;
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
    hasRole(role: Role): boolean {
      if (this.hasAdmin) {
        return true;
      }
      return this.user.role >= role;
    },
    hasRoleForGroup(role: Role, group: FullGroup): boolean {
      if (this.hasAdmin) {
        return true;
      }
      if (this.user.role >= role) {
        return true;
      }
      return group.members.some((member) => member.role >= role && member.id === this.user.id);
    }
  },
  getters: {
    hasRoot(): boolean {
      return this.user.isRoot;
    },
    hasAdmin(): boolean {
      return this.hasRoot || this.user.isAdmin;
    }
  }
});