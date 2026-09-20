export function formatHRDate(iso: string): string {
    const date = new Date(iso);
    const now = new Date();

    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const startOfYesterday = new Date(startOfToday);
    startOfYesterday.setDate(startOfYesterday.getDate() - 1);

    const startOfWeek = new Date(startOfToday);
    startOfWeek.setDate(startOfWeek.getDate() - 7);

    if (date >= startOfToday) {
        return date.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        });
    }

    if (date >= startOfYesterday) {
        return "Yesterday";
    }

    if (date >= startOfWeek) {
        const days = Math.floor(
            (startOfToday.getTime() - date.getTime()) / (1000 * 60 * 60 * 24),
        );

        return `${days} days ago`;
    }

    return date.toLocaleDateString("en-GB", {
        day: "numeric",
        month: "long",
    });
}
