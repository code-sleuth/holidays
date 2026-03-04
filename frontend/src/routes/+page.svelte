<script lang="ts">
    import { onMount } from "svelte";
    import {
        Card,
        CardHeader,
        CardTitle,
        CardContent,
    } from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import {
        loadHolidays,
        loadLookups,
        getHolidays,
        getLookups,
        calculateDays,
    } from "$lib/stores/holiday.svelte";

    let selectedValue = $state("");

    let holidays = $derived(getHolidays());
    let lookups = $derived(getLookups());

    async function handleCalculate() {
        if (!selectedValue) return;
        const [date, ...nameParts] = selectedValue.split("|");
        const name = nameParts.join("|");
        await calculateDays(name, date);
    }

    onMount(() => {
        loadHolidays();
        loadLookups();
    });
</script>

<div class="min-h-screen bg-slate-50 p-8">
    <div class="max-w-2xl mx-auto">
        <h1 class="text-3xl font-bold text-slate-900 mb-8">
            Days Until Holiday
        </h1>

        <!-- Calculator -->
        <Card class="mb-8">
            <CardHeader>
                <CardTitle>Select a Holiday</CardTitle>
            </CardHeader>
            <CardContent>
                <div class="flex gap-3">
                    <select
                        bind:value={selectedValue}
                        class="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400"
                    >
                        <option value="">Choose a holiday...</option>
                        {#each holidays as holiday}
                            <option value="{holiday.date}|{holiday.localName}">
                                {holiday.localName} ({holiday.date})
                            </option>
                        {/each}
                    </select>
                    <Button onclick={handleCalculate} disabled={!selectedValue}
                        >Calculate</Button
                    >
                </div>
            </CardContent>
        </Card>

        <!-- History -->
        <h2 class="text-xl font-semibold text-slate-900 mb-4">
            Lookup History
        </h2>

        {#if lookups.length === 0}
            <Card>
                <CardContent class="py-8 text-center text-slate-500">
                    No lookups yet. Select a holiday and click Calculate.
                </CardContent>
            </Card>
        {:else}
            <div class="space-y-3">
                {#each lookups as lookup (lookup.id)}
                    <Card>
                        <CardContent
                            class="py-4 flex items-center justify-between"
                        >
                            <div>
                                <p class="font-medium text-slate-900">
                                    {lookup.holidayName}
                                </p>
                                <p class="text-sm text-slate-500">
                                    {lookup.holidayDate}
                                </p>
                            </div>
                            <div class="text-right">
                                <p
                                    class="text-2xl font-bold {lookup.daysUntil >
                                    0
                                        ? 'text-green-600'
                                        : lookup.daysUntil < 0
                                          ? 'text-red-500'
                                          : 'text-slate-900'}"
                                >
                                    {lookup.daysUntil}
                                </p>
                                <p class="text-sm text-slate-500">
                                    {lookup.daysUntil === 1 ? "day" : "days"}
                                    {lookup.daysUntil >= 0 ? "until" : "ago"}
                                </p>
                            </div>
                        </CardContent>
                    </Card>
                {/each}
            </div>
        {/if}
    </div>
</div>
