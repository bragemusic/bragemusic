export default function TightLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <main className="overflow-y-scroll flex-grow mx-auto cursor-default select-none max-h-max">
      {children}
    </main>
  );
}
