<template>
  <div />
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'

@Component
export default class PreloadImage extends Vue {
  @Prop()
  public src!: string

  @Prop({ default: () => 200 })
  public delayMs!: number

  public mounted() {
    this.delayedPreloadImage()
  }

  @Watch('src')
  public onSrcChanged() {
    this.delayedPreloadImage()
  }

  private delayedPreloadImage() {
    setTimeout(() => {
      this.preloadImage()
    }, this.delayMs)
  }

  private preloadImage() {
    const img = new Image()
    img.src = this.src
  }
}
</script>
