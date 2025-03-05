import type { Knex } from "knex";

const TABLE_PAYMENT = 'payments'

export async function up(knex: Knex): Promise<void> {
    if (!(await knex.schema.hasTable(TABLE_PAYMENT))) {
        await knex.schema.createTable("payments", (table) => {
            table.uuid("id", { primaryKey: true, useBinaryUuid: true }).defaultTo(knex.fn.uuid());;
            table.uuid("order_id").notNullable();
            table.string("onl_payment_intent_id", 256).nullable().unique();
            table.integer("amount").notNullable();
            table.enu('status', ['pending', 'success', 'cancel']).defaultTo('pending');
            table.enu('payment_method', ['cash', 'credit card']).notNullable();

            table.timestamp('created_at').defaultTo(knex.fn.now());
            table.timestamp('updated_at').defaultTo(knex.fn.now());
        })
    }
}

export async function down(knex: Knex): Promise<void> {
    await knex.schema.dropTableIfExists(TABLE_PAYMENT)
}
