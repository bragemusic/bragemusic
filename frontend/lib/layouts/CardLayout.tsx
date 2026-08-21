import PageTitleLayout from "./page_title";

type CardLayoutProps = {
  title?: string;
  headerContent?: React.ReactNode;
  children: React.ReactNode;
};

export default function CardLayout({
  title,
  headerContent,
  children,
}: CardLayoutProps) {
  return (
    <PageTitleLayout headerContent={headerContent} title={title}>
      <div />
      <div className="grid flex-shrink gap-2 w-full sm:gap-6 sm:w-auto md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-6">
        {children}
      </div>
      <div />
    </PageTitleLayout>
  );
}
