import { ReactNode } from 'react'
import { Sidebar } from './Sidebar'
import { Header } from './Header'

export function PageLayout({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="flex h-screen bg-slate-50">
      <Sidebar />
      <div className="flex-1 flex flex-col min-w-0">
        <Header title={title} />
        <main className="flex-1 overflow-y-auto p-6">
          {children}
        </main>
      </div>
    </div>
  )
}
