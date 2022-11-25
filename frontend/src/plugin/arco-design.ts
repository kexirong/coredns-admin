import { Space, Form, Button, Input, Card,Layout } from '@arco-design/web-vue'
import { FormItem, CardMeta,LayoutContent,LayoutHeader,LayoutSider } from '@arco-design/web-vue/es'
import { App } from 'vue'
export default (app: App<Element>) => {
    app.component('AForm', Form)
    app.component('ASpace', Space)
    app.component('AButton', Button)
    app.component('AInput', Input)
    app.component('AFormItem', FormItem)
    app.component('ACard', Card)
    app.component('ACardMeta', CardMeta)
    app.component('ALayout', Layout)
    app.component('ALayoutContent', LayoutContent)
    app.component('ALayoutHeader', LayoutHeader)
    app.component('ALayoutSider', LayoutSider)
}
