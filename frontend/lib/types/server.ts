
export namespace responses {
	export class ListPaginationPayload<T> {
	    items?: T[];
	    page: number;
	    limit: number;
	    total_pages: number;
	    total_items: number;

      static createFrom<T>(source: any = {}, itemClass?: any): ListPaginationPayload<T> {
	        return new ListPaginationPayload<T>(source, itemClass);
	    }

      constructor(source: any = {}, itemClass?: any) {
	        if ('string' === typeof source) source = JSON.parse(source);

          this.items = this.convertValues(source["items"], itemClass);
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.total_pages = source["total_pages"];
	        this.total_items = source["total_items"];
	    }

    private convertValues(a: any, classs?: any): any {
      if (!a) return a;

      if (Array.isArray(a)) {
        return a.map((elem) =>
          classs ? new classs(elem) : elem
        );
      }

      return classs ? new classs(a) : a;
    }
  }

 export class ListPayload<T> {
    items?: T[];
    count: number;

    static createFrom<T>(source: any = {}, itemClass?: any): ListPayload<T> {
      return new ListPayload<T>(source, itemClass);
    }

    constructor(source: any = {}, itemClass?: any) {
      if (typeof source === "string") source = JSON.parse(source);

      this.count = source["count"];
      this.items = this.convertValues(source["items"], itemClass);
    }

    private convertValues(a: any, classs?: any): any {
      if (!a) return a;

      if (Array.isArray(a)) {
        return a.map((elem) =>
          classs ? new classs(elem) : elem
        );
      }

      return classs ? new classs(a) : a;
    }
  }

  export class Login {
    token: string;
    token_type: string;
    expires_in: number;

    static createFrom(source: any = {}) {
      return new Login(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.token = source["token"];
      this.token_type = source["token_type"];
      this.expires_in = source["expires_in"];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice && a.map) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ("object" === typeof a) {
        if (asMap) {
          for (const key of Object.keys(a)) {
            a[key] = new classs(a[key]);
          }

          return a;
        }

        return new classs(a);
      }

      return a;
    }
  }

  export class DeviceRegister {
    id: string;

    static createFrom(source: any = {}) {
      return new DeviceRegister(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.id= source["id"];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice && a.map) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ("object" === typeof a) {
        if (asMap) {
          for (const key of Object.keys(a)) {
            a[key] = new classs(a[key]);
          }

          return a;
        }

        return new classs(a);
      }

      return a;
    }
  }

  export class CreateApiToken {
    token: string;

    static createFrom(source: any = {}) {
      return new CreateApiToken(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.token= source["token"];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice && a.map) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ("object" === typeof a) {
        if (asMap) {
          for (const key of Object.keys(a)) {
            a[key] = new classs(a[key]);
          }

          return a;
        }

        return new classs(a);
      }

      return a;
    }
  }

}

export namespace requests {

  export class DevicesRegister {
      id?: string;
      name: string;
      type: string;
      interface: string;
      icon: string;
      supports_playback: boolean;
      platform: string;
      version: string;

      static createFrom(source: any = {}) {
          return new DevicesRegister(source);
      }

      constructor(source: any = {}) {
          if ('string' === typeof source) source = JSON.parse(source);
          this.id = source["id"];
          this.name = source["name"];
          this.type = source["type"];
          this.interface = source["interface"];
          this.icon = source["icon"];
          this.supports_playback = source["supports_playback"];
          this.platform = source["platform"];
          this.version = source["version"];
      }
  }

}
