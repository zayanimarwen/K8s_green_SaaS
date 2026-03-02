interface Props { isLive: boolean }
export function LiveIndicator({ isLive }: Props) {
  return (
    <div className="flex items-center gap-1.5 text-xs">
      <span className={`w-2 h-2 rounded-full ${isLive ? 'bg-green-500 animate-pulse' : 'bg-slate-300'}`} />
      <span className={isLive ? 'text-green-600 font-medium' : 'text-slate-400'}>
        {isLive ? 'Live' : 'Polling'}
      </span>
    </div>
  )
}
