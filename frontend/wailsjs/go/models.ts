export namespace models {
	
	export class Account {
	    id: number;
	    name: string;
	    type: string;
	    balance: number;
	    currency: string;
	    color: string;
	    icon: string;
	    description: string;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.balance = source["balance"];
	        this.currency = source["currency"];
	        this.color = source["color"];
	        this.icon = source["icon"];
	        this.description = source["description"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class AccountBalance {
	    account: Account;
	    income: number;
	    expense: number;
	
	    static createFrom(source: any = {}) {
	        return new AccountBalance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account = this.convertValues(source["account"], Account);
	        this.income = source["income"];
	        this.expense = source["expense"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class KanbanColumn {
	    id: number;
	    board_id: number;
	    name: string;
	    color: string;
	    position: number;
	    wip_limit?: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new KanbanColumn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.board_id = source["board_id"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.position = source["position"];
	        this.wip_limit = source["wip_limit"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class BoardWithColumns {
	    id: number;
	    name: string;
	    profile_id: number;
	    is_shared: boolean;
	    // Go type: time
	    created_at: any;
	    columns: KanbanColumn[];
	
	    static createFrom(source: any = {}) {
	        return new BoardWithColumns(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.profile_id = source["profile_id"];
	        this.is_shared = source["is_shared"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.columns = this.convertValues(source["columns"], KanbanColumn);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class Category {
	    id: number;
	    name: string;
	    color: string;
	    type: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.type = source["type"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class CreateProfileInput {
	    name: string;
	    email: string;
	    role: string;
	    avatar_url: string;
	    pin?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProfileInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.role = source["role"];
	        this.avatar_url = source["avatar_url"];
	        this.pin = source["pin"];
	    }
	}
	export class Event {
	    id: number;
	    title: string;
	    description: string;
	    // Go type: time
	    start_date: any;
	    // Go type: time
	    end_date: any;
	    all_day: boolean;
	    color: string;
	    location: string;
	    reminder_minutes?: number;
	    google_event_id?: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.end_date = this.convertValues(source["end_date"], null);
	        this.all_day = source["all_day"];
	        this.color = source["color"];
	        this.location = source["location"];
	        this.reminder_minutes = source["reminder_minutes"];
	        this.google_event_id = source["google_event_id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class FinanceFilter {
	    account_id?: number;
	    type: string;
	    category: string;
	    date_from: string;
	    date_to: string;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new FinanceFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account_id = source["account_id"];
	        this.type = source["type"];
	        this.category = source["category"];
	        this.date_from = source["date_from"];
	        this.date_to = source["date_to"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class Transaction {
	    id: number;
	    account_id: number;
	    account_name?: string;
	    to_account_id?: number;
	    to_account_name?: string;
	    type: string;
	    category: string;
	    description: string;
	    amount: number;
	    // Go type: time
	    date: any;
	    is_recurring: boolean;
	    recurrence_rule: string;
	    tags: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Transaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.account_id = source["account_id"];
	        this.account_name = source["account_name"];
	        this.to_account_id = source["to_account_id"];
	        this.to_account_name = source["to_account_name"];
	        this.type = source["type"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.amount = source["amount"];
	        this.date = this.convertValues(source["date"], null);
	        this.is_recurring = source["is_recurring"];
	        this.recurrence_rule = source["recurrence_rule"];
	        this.tags = source["tags"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class FinanceSummary {
	    total_balance: number;
	    total_income: number;
	    total_expense: number;
	    accounts: AccountBalance[];
	    recent_transactions: Transaction[];
	
	    static createFrom(source: any = {}) {
	        return new FinanceSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_balance = source["total_balance"];
	        this.total_income = source["total_income"];
	        this.total_expense = source["total_expense"];
	        this.accounts = this.convertValues(source["accounts"], AccountBalance);
	        this.recent_transactions = this.convertValues(source["recent_transactions"], Transaction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	
	export class Label {
	    id: number;
	    name: string;
	    color: string;
	    board_id: number;
	
	    static createFrom(source: any = {}) {
	        return new Label(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.board_id = source["board_id"];
	    }
	}
	export class Note {
	    id: number;
	    title: string;
	    content: string;
	    category_id?: number;
	    is_favorite: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Note(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.category_id = source["category_id"];
	        this.is_favorite = source["is_favorite"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class Product {
	    id: number;
	    name: string;
	    category: string;
	    unit: string;
	    min_quantity: number;
	    current_stock: number;
	    price: number;
	    notes: string;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.unit = source["unit"];
	        this.min_quantity = source["min_quantity"];
	        this.current_stock = source["current_stock"];
	        this.price = source["price"];
	        this.notes = source["notes"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class ProfileSummary {
	    id: number;
	    name: string;
	    email: string;
	    role: string;
	    avatar_url: string;
	    google_calendar_id: string;
	    has_google_token: boolean;
	    has_pin: boolean;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ProfileSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.role = source["role"];
	        this.avatar_url = source["avatar_url"];
	        this.google_calendar_id = source["google_calendar_id"];
	        this.has_google_token = source["has_google_token"];
	        this.has_pin = source["has_pin"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class ScheduledTransaction {
	    id: number;
	    account_id: number;
	    account_name?: string;
	    type: string;
	    category: string;
	    description: string;
	    amount: number;
	    // Go type: time
	    scheduled_at: any;
	    is_recurring: boolean;
	    recurrence_rule: string;
	    is_executed: boolean;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ScheduledTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.account_id = source["account_id"];
	        this.account_name = source["account_name"];
	        this.type = source["type"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.amount = source["amount"];
	        this.scheduled_at = this.convertValues(source["scheduled_at"], null);
	        this.is_recurring = source["is_recurring"];
	        this.recurrence_rule = source["recurrence_rule"];
	        this.is_executed = source["is_executed"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class ShoppingItem {
	    id: number;
	    shopping_list_id: number;
	    product_id?: number;
	    name: string;
	    quantity: number;
	    unit: string;
	    estimated_price: number;
	    actual_price: number;
	    category: string;
	    is_bought: boolean;
	    notes: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ShoppingItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.shopping_list_id = source["shopping_list_id"];
	        this.product_id = source["product_id"];
	        this.name = source["name"];
	        this.quantity = source["quantity"];
	        this.unit = source["unit"];
	        this.estimated_price = source["estimated_price"];
	        this.actual_price = source["actual_price"];
	        this.category = source["category"];
	        this.is_bought = source["is_bought"];
	        this.notes = source["notes"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class ShoppingList {
	    id: number;
	    name: string;
	    month: string;
	    description: string;
	    is_completed: boolean;
	    total_budget: number;
	    total_spent: number;
	    items?: ShoppingItem[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ShoppingList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.month = source["month"];
	        this.description = source["description"];
	        this.is_completed = source["is_completed"];
	        this.total_budget = source["total_budget"];
	        this.total_spent = source["total_spent"];
	        this.items = this.convertValues(source["items"], ShoppingItem);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class StockAlert {
	    product: Product;
	    deficit: number;
	
	    static createFrom(source: any = {}) {
	        return new StockAlert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product = this.convertValues(source["product"], Product);
	        this.deficit = source["deficit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class StockMovement {
	    id: number;
	    product_id: number;
	    product_name?: string;
	    type: string;
	    quantity: number;
	    price: number;
	    notes: string;
	    // Go type: time
	    date: any;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new StockMovement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.type = source["type"];
	        this.quantity = source["quantity"];
	        this.price = source["price"];
	        this.notes = source["notes"];
	        this.date = this.convertValues(source["date"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class Subtask {
	    id: number;
	    task_id: number;
	    title: string;
	    completed: boolean;
	    position: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Subtask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task_id = source["task_id"];
	        this.title = source["title"];
	        this.completed = source["completed"];
	        this.position = source["position"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class Task {
	    id: number;
	    title: string;
	    description: string;
	    status: string;
	    priority: string;
	    category_id?: number;
	    board_id?: number;
	    column_id?: number;
	    position: number;
	    // Go type: time
	    due_date?: any;
	    // Go type: time
	    start_date?: any;
	    estimated_hours?: number;
	    progress: number;
	    // Go type: time
	    completed_at?: any;
	    google_event_id?: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.category_id = source["category_id"];
	        this.board_id = source["board_id"];
	        this.column_id = source["column_id"];
	        this.position = source["position"];
	        this.due_date = this.convertValues(source["due_date"], null);
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.estimated_hours = source["estimated_hours"];
	        this.progress = source["progress"];
	        this.completed_at = this.convertValues(source["completed_at"], null);
	        this.google_event_id = source["google_event_id"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class TaskActivity {
	    id: number;
	    task_id: number;
	    profile_id?: number;
	    profile?: ProfileSummary;
	    action: string;
	    from_value: string;
	    to_value: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TaskActivity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task_id = source["task_id"];
	        this.profile_id = source["profile_id"];
	        this.profile = this.convertValues(source["profile"], ProfileSummary);
	        this.action = source["action"];
	        this.from_value = source["from_value"];
	        this.to_value = source["to_value"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class TaskComment {
	    id: number;
	    task_id: number;
	    profile_id: number;
	    profile?: ProfileSummary;
	    content: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TaskComment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task_id = source["task_id"];
	        this.profile_id = source["profile_id"];
	        this.profile = this.convertValues(source["profile"], ProfileSummary);
	        this.content = source["content"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class TaskDetail {
	    assignees: ProfileSummary[];
	    subtasks: Subtask[];
	    labels: Label[];
	    comments: TaskComment[];
	    dependencies: number[];
	    dependents: number[];
	
	    static createFrom(source: any = {}) {
	        return new TaskDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assignees = this.convertValues(source["assignees"], ProfileSummary);
	        this.subtasks = this.convertValues(source["subtasks"], Subtask);
	        this.labels = this.convertValues(source["labels"], Label);
	        this.comments = this.convertValues(source["comments"], TaskComment);
	        this.dependencies = source["dependencies"];
	        this.dependents = source["dependents"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class TaskFilter {
	    Status: string;
	    Priority: string;
	    CategoryID?: number;
	    BoardID?: number;
	    ColumnID?: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Status = source["Status"];
	        this.Priority = source["Priority"];
	        this.CategoryID = source["CategoryID"];
	        this.BoardID = source["BoardID"];
	        this.ColumnID = source["ColumnID"];
	    }
	}
	
	export class WAChat {
	    jid: string;
	    name: string;
	    last_message: string;
	    // Go type: time
	    last_message_at: any;
	    unread_count: number;
	
	    static createFrom(source: any = {}) {
	        return new WAChat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.jid = source["jid"];
	        this.name = source["name"];
	        this.last_message = source["last_message"];
	        this.last_message_at = this.convertValues(source["last_message_at"], null);
	        this.unread_count = source["unread_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class WAMessage {
	    id: string;
	    chat_jid: string;
	    sender_jid: string;
	    content: string;
	    is_from_me: boolean;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new WAMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.chat_jid = source["chat_jid"];
	        this.sender_jid = source["sender_jid"];
	        this.content = source["content"];
	        this.is_from_me = source["is_from_me"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class WAScheduleInput {
	    chat_jid: string;
	    chat_name: string;
	    content: string;
	    // Go type: time
	    scheduled_at: any;
	
	    static createFrom(source: any = {}) {
	        return new WAScheduleInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chat_jid = source["chat_jid"];
	        this.chat_name = source["chat_name"];
	        this.content = source["content"];
	        this.scheduled_at = this.convertValues(source["scheduled_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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
	export class WAScheduled {
	    id: number;
	    chat_jid: string;
	    chat_name: string;
	    content: string;
	    // Go type: time
	    scheduled_at: any;
	    sent: boolean;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new WAScheduled(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.chat_jid = source["chat_jid"];
	        this.chat_name = source["chat_name"];
	        this.content = source["content"];
	        this.scheduled_at = this.convertValues(source["scheduled_at"], null);
	        this.sent = source["sent"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
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

export namespace services {
	
	export class WAStatusInfo {
	    status: string;
	    phone: string;
	
	    static createFrom(source: any = {}) {
	        return new WAStatusInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.phone = source["phone"];
	    }
	}

}

